INSERT INTO api.finite_automatas (id, description, tuple, render) VALUES
  (
    '11111111-1111-1111-1111-111111111111',
    'Simple DFA with two states',
    '{
      "alphabet": ["a", "b"],
      "states": ["q0", "q1"],
      "initial": "q0",
      "acceptance": ["q1"],
      "transitions": [
        ["q1", "q0"],
        ["q0", "q1"]
      ]
    }',
    '/data/images/11111111-1111-1111-1111-111111111111.svg'
  ),
  (
    '22222222-2222-2222-2222-222222222222',
    'NFA that accepts strings ending in "ab"',
    '{
      "alphabet": ["a", "b"],
      "states": ["s0", "s1", "s2"],
      "initial": "s0",
      "acceptance": ["s2"],
      "transitions": [
        [["s0", "s1"], "@v"],
        ["@v", "s2"],
        ["@v", "@v"]
      ]
    }',
    '/data/images/22222222-2222-2222-2222-222222222222.svg'
  );
