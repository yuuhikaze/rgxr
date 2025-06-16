INSERT INTO api.finite_automatas (id, description, tuple, render) VALUES
  (
    'cf6d1645-0856-46ed-ae20-2e0852ec9b2b',
    '#1.3',
    '{
      "alphabet": ["u", "d"],
      "states": ["q1", "q2", "q3", "q4", "q5"],
      "initial": "q3",
      "acceptance": ["q3"],
      "transitions": [
        ["q1", "q2"],
        ["q1", "q3"],
        ["q2", "q4"],
        ["q3", "q5"],
        ["q4", "q5"]
      ]
    }',
    '/path/to/render'
  );
