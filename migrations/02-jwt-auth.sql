CREATE EXTENSION pgcrypto;

-- https://github.com/michelp/pgjwt/blob/master/pgjwt--0.2.0.sql
CREATE OR REPLACE FUNCTION url_encode (data bytea)
    RETURNS text
    LANGUAGE sql
    AS $$
    SELECT
        translate(encode(data, 'base64'), E'+/=\n', '-_');
$$ IMMUTABLE;

CREATE OR REPLACE FUNCTION url_decode (data text)
    RETURNS bytea
    LANGUAGE sql
    AS $$
    WITH t AS (
        SELECT
            translate(data, '-_', '+/') AS trans
),
rem AS (
    SELECT
        length(t.trans) % 4 AS remainder
    FROM
        t) -- compute padding size
    SELECT
        decode(t.trans || CASE WHEN rem.remainder > 0 THEN
                repeat('=', (4 - rem.remainder))
            ELSE
                ''
            END, 'base64')
    FROM
        t,
        rem;
$$ IMMUTABLE;

CREATE OR REPLACE FUNCTION algorithm_sign (signables text, secret text, algorithm text)
    RETURNS text
    LANGUAGE sql
    AS $$
    WITH alg AS (
        SELECT
            CASE WHEN algorithm = 'HS256' THEN
                'sha256'
            WHEN algorithm = 'HS384' THEN
                'sha384'
            WHEN algorithm = 'HS512' THEN
                'sha512'
            ELSE
                ''
            END AS id) -- hmac throws error
        SELECT
            url_encode (hmac(signables, secret, alg.id))
        FROM
            alg;
$$ IMMUTABLE;

CREATE OR REPLACE FUNCTION sign (payload json, secret text, algorithm text DEFAULT 'HS256')
    RETURNS text
    LANGUAGE sql
    AS $$
    WITH header AS (
        SELECT
            url_encode (convert_to('{"alg":"' || algorithm || '","typ":"JWT"}', 'utf8')) AS data
),
payload AS (
    SELECT
        url_encode (convert_to(payload::text, 'utf8')) AS data
),
signables AS (
    SELECT
        header.data || '.' || payload.data AS data
    FROM
        header,
        payload
)
SELECT
    signables.data || '.' || algorithm_sign (signables.data, secret, algorithm)
FROM
    signables;
$$ IMMUTABLE;

CREATE OR REPLACE FUNCTION try_cast_double (inp text)
    RETURNS double precision
    AS $$
BEGIN
    BEGIN
        RETURN inp::double precision;
    EXCEPTION
        WHEN OTHERS THEN
            RETURN NULL;
    END;
END;

$$
LANGUAGE plpgsql
IMMUTABLE;

CREATE OR REPLACE FUNCTION verify (token text, secret text, algorithm text DEFAULT 'HS256')
    RETURNS TABLE (
        header json,
        payload json,
        valid boolean)
    LANGUAGE sql
    AS $$
    SELECT
        jwt.header AS header,
        jwt.payload AS payload,
        jwt.signature_ok
        AND tstzrange(to_timestamp(try_cast_double (jwt.payload ->> 'nbf')), to_timestamp(try_cast_double (jwt.payload ->> 'exp'))) @> CURRENT_TIMESTAMP AS valid
    FROM (
        SELECT
            convert_from(url_decode (r[1]), 'utf8')::json AS header,
            convert_from(url_decode (r[2]), 'utf8')::json AS payload,
            r[3] = algorithm_sign (r[1] || '.' || r[2], secret, algorithm) AS signature_ok
        FROM
            regexp_split_to_array(token, '\.') r) jwt
$$ IMMUTABLE;

-- https://docs.postgrest.org/en/v13/how-tos/sql-user-management.html
-- change this in the future
ALTER DATABASE postgres SET "app.jwt_secret" TO 'reallyreallyreallyreallyverysafe';

CREATE SCHEMA IF NOT EXISTS basic_auth;

CREATE TABLE basic_auth.users (
    email text PRIMARY KEY CHECK (email ~* '^.+@.+\..+$'),
    pass text NOT NULL CHECK (length(pass) < 512),
    role name NOT NULL CHECK (length(ROLE) < 512)
);

CREATE FUNCTION basic_auth.check_role_exists ()
    RETURNS TRIGGER
    AS $$
BEGIN
    IF NOT EXISTS (
        SELECT
            1
        FROM
            pg_roles AS r
        WHERE
            r.rolname = NEW.role) THEN
    RAISE foreign_key_violation
    USING message = 'unknown database role: ' || NEW.role;
    RETURN NULL;
END IF;
    RETURN new;
END
$$
LANGUAGE plpgsql;

CREATE CONSTRAINT TRIGGER ensure_user_role_exists
    AFTER INSERT OR UPDATE ON basic_auth.users
    FOR EACH ROW
    EXECUTE PROCEDURE basic_auth.check_role_exists ();

CREATE FUNCTION basic_auth.encrypt_pass ()
    RETURNS TRIGGER
    AS $$
BEGIN
    IF tg_op = 'INSERT' OR NEW.pass <> OLD.pass THEN
        NEW.pass = crypt(NEW.pass, gen_salt('bf'));
    END IF;
    RETURN new;
END
$$
LANGUAGE plpgsql;

CREATE TRIGGER encrypt_pass
    BEFORE INSERT OR UPDATE ON basic_auth.users
    FOR EACH ROW
    EXECUTE PROCEDURE basic_auth.encrypt_pass ();

CREATE FUNCTION basic_auth.user_role (email text, pass text)
    RETURNS name
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN (
        SELECT
            ROLE
        FROM
            basic_auth.users
        WHERE
            users.email = user_role.email
            AND users.pass = crypt(user_role.pass, users.pass));
END;
$$;

CREATE FUNCTION api.login (email text, pass text, out token text)
AS $$
DECLARE
    _role name;
BEGIN
    -- check email and password
    SELECT
        basic_auth.user_role (email, pass) INTO _role;
    IF _role IS NULL THEN
        RAISE invalid_password
        USING message = 'invalid user or password';
    END IF;
        SELECT
            sign(row_to_json(r), current_setting('app.jwt_secret')) AS token
        FROM (
            SELECT
                _role AS role,
                login.email AS email,
                extract(epoch FROM now())::integer + 60 * 60 AS exp) r INTO token;
END;
$$
LANGUAGE plpgsql
SECURITY DEFINER;

GRANT EXECUTE ON FUNCTION api.login (text, text) TO web_anon;

-- web_editor role
CREATE ROLE web_editor noinherit;

GRANT SELECT, INSERT, UPDATE, DELETE ON api.finite_automatas TO web_editor;

GRANT web_editor TO authenticator;

INSERT INTO basic_auth.users (email, pass, role)
    VALUES ('editor@example.com', 'securepassword', 'web_editor');
