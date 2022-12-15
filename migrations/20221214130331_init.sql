-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS slot (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description TEXT
);

CREATE TABLE IF NOT EXISTS banner (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description TEXT
);

CREATE TABLE IF NOT EXISTS social_group (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description TEXT
);

CREATE TABLE IF NOT EXISTS stat (
    banner_id UUID,
    slot_id   UUID,
    social_group_id  UUID,
    shows     INT,
    clicks    INT,
    PRIMARY KEY (banner_id, slot_id, social_group_id),
    FOREIGN KEY (banner_id) REFERENCES banner (id),
    FOREIGN KEY (slot_id) REFERENCES slot (id),
    FOREIGN KEY (social_group_id) REFERENCES social_group (id)
);

CREATE TABLE IF NOT EXISTS banner_slot(
    banner_id UUID,
    slot_id UUID,
    PRIMARY KEY (banner_id, slot_id),
    FOREIGN KEY (banner_id) REFERENCES banner (id),
    FOREIGN KEY (slot_id) REFERENCES slot (id)
);

INSERT INTO slot(id, description)
VALUES ('00000000-0000-0000-0000-000000000001', 'slot 1'),
       ('00000000-0000-0000-0000-000000000002', 'slot 2'),
       ('00000000-0000-0000-0000-000000000003', 'slot 3'),
       ('00000000-0000-0000-0000-000000000004', 'slot 4'),
       ('00000000-0000-0000-0000-000000000005', 'slot 5'),
       ('00000000-0000-0000-0000-000000000006', 'slot 6');

INSERT INTO banner(id, description)
VALUES ('00000000-0000-0000-1111-000000000001', 'banner 1'),
       ('00000000-0000-0000-1111-000000000002', 'banner 2'),
       ('00000000-0000-0000-1111-000000000003', 'banner 3'),
       ('00000000-0000-0000-1111-000000000004', 'banner 4'),
       ('00000000-0000-0000-1111-000000000005', 'banner 5');

INSERT INTO social_group(id, description)
VALUES ('00000000-0000-0000-2222-000000000001', 'social group 1'),
       ('00000000-0000-0000-2222-000000000002', 'social group 2'),
       ('00000000-0000-0000-2222-000000000003', 'social group 3'),
       ('00000000-0000-0000-2222-000000000004', 'social group 4'),
       ('00000000-0000-0000-2222-000000000005', 'social group 5');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stat;
DROP TABLE IF EXISTS banner;
DROP TABLE IF EXISTS slot;
DROP TABLE IF EXISTS social_group;
-- +goose StatementEnd