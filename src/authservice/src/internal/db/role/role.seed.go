package role

var SQL = `
insert into public.roles (id, updated_at, created_at, deleted_at, name, code)
values  ('39a5635c-be8a-4ccd-b7a6-35b2c9f58819', 11, 1, 0, 'admin', 2),
        ('aace9864-752f-4959-ad22-5a8263046bff', 11, 1, 0, 'customer', 0),
        ('bb2438b6-256b-48b0-b1ad-807324fed0cc', 11, 1, 0, 'handyman', 1);
`
