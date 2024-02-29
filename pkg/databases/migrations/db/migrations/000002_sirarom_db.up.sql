BEGIN;

INSERT INTO
    "roles" ("title")
VALUES ('dev'),
    ('admin');

INSERT INTO
    "users" (
        "username", "name", "email", "password", "role_id"
    )
VALUES (
        'y.pornwisa','pornwisa' ,'y.pornwisa@gmail.com', '$2a$10$8KzaNdKIMyOkASCH4QvSKuEMIY7Jc3vcHDuSJvXLii1rvBNgz60a6', 1
    ),
    (
        'admin', 'admin', 'admin001@kawaii.com', '$2a$10$3qqNPE.TJpNGYCohjTgw9.v1z0ckovx95AmiEtUXcixGAgfW7.wCi', 2
    );

COMMIT;