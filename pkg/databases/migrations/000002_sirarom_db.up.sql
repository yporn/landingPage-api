BEGIN;

INSERT INTO "roles" ("title") VALUES ('all'), ('home'), ('project'), ('promotion'), ('activity'), ('job');

INSERT INTO
    "users" (
        "username", "name", "email", "password", "tel", "display"
    )
VALUES (
        'y.pornwisa', 
        'pornwisa', 
        'y.pornwisa@gmail.com', 
        '$2a$10$8KzaNdKIMyOkASCH4QvSKuEMIY7Jc3vcHDuSJvXLii1rvBNgz60a6', 
        '0900000000',
        'published'
    ),
    (
        'admin', 
        'admin', 
        'admin001@kawaii.com', 
        '$2a$10$3qqNPE.TJpNGYCohjTgw9.v1z0ckovx95AmiEtUXcixGAgfW7.wCi', 
        '0900000000',
        'published'
    );

INSERT INTO "user_roles" ("user_id", "role_id") VALUES (2, 2), (2, 3), (2, 4), (1, 1);

INSERT INTO banners ("index", "delay", "display")
VALUES (1, 5, 'published'),
       (2, 7, 'unpublished');

INSERT INTO
    "careers" (
        "position", "amount", "location", "description", 
        "qualification", "start_date", "end_date", "status", "display"
    )
VALUES (
        'พนักงานเสิร์ฟชาย', '2', 
        'ประจำบ้านพักผู้บริหาร ต.จอหอ อ.เมืองนครราชสีมา จ.นครราชสีมา', 
        '"มีหน้าที่ต้อนรับและเสิร์ฟอาหารผู้บริหาร ดูแลรักษามาตรฐานความปลอดภัยและสุขอนามัยของอาหาร"', 
        'เพศชาย อายุ 25-35 ปี', 
        '2024-05-03', '2024-05-15', 'opening', 'published'
    ),
    (
        'วิศวกรประจำโครงการ', '1', 
        'อ.เมืองนครราชสีมา จ.นครราชสีมา', 
        'จัดเตรียมความพร้อมและวางแผนการทำงาน วางแผนการก่อสร้างโดยประสานและสอดคล้อง', 
        'เพศชาย อายุ 30 ปีขึ้นไป วุฒิปริญญาตรี สาขาวิศวกรรมโยธา', 
        '2024-02-11', '2024-02-18', 'closed', 'unpublished'
    ),
    (
        'สถาปนิกโครงการ', '2', 'อ.เมืองนครราชสีมา จ.นครราชสีมา', 
        'ออกแบบงานด้านสถาปัตยกรรมให้ตรงตาม แผนงานที่กำหนดไว้ตรวจสอบความถูกต้องของแบบและรายการ', 
        'ไม่จำกัดเพศ อาย 25-35 ปี', 
        '2024-03-02', '2024-03-30', 'opening', 'published'
    );

INSERT INTO projects (name, "index", heading, text, location, price, status_project, type_project, description, name_facebook, link_facebook, tel, address, link_location, display)
VALUES ('Project 1', 1, 'Project 1 Heading', 'Project 1 Text', 'Project 1 Location', 100000, 'ready', 'present', 'Project 1 Description', 'Project 1 Facebook', 'https://facebook.com/project1', '123456789', 'Project 1 Address', 'https://maps.google.com/project1', 'published'),
       ('Project 2', 2, 'Project 2 Heading', 'Project 2 Text', 'Project 2 Location', 200000, 'new', 'future', 'Project 2 Description', 'Project 2 Facebook', 'https://facebook.com/project2', '987654321', 'Project 2 Address', 'https://maps.google.com/project2', 'published');

-- Seed data for the "project_house_type_items" table
INSERT INTO project_house_type_items (project_id, name)
VALUES (1, 'Type A'),
       (1, 'Type B'),
       (2, 'Type C');

-- Seed data for the "project_desc_area_items" table
INSERT INTO project_desc_area_items (project_id, item, amount, unit)
VALUES (1, 'Area A', '100', 'sqm'),
       (1, 'Area B', '200', 'sqm'),
       (2, 'Area C', '150', 'sqm');

-- Seed data for the "project_comfortable_items" table
INSERT INTO project_comfortable_items (project_id, item)
VALUES (1, 'Comfortable A'),
       (1, 'Comfortable B'),
       (2, 'Comfortable C');


INSERT INTO "house_models" 
	("project_id", "index", "name", "description", "link_video", "link_virtual_tour", "display")
VALUES 
	(1, 1, 'เติมรัก', 'เติมเต็มความสุข.. กับการเริ่มต้นของครอบครัว', 'https://www.youtube.com/watch?v=BxuY9FET9Y4&list=RDHTcL9WkB_wg&index=15', 'test', 'published');
	
INSERT INTO "house_model_type_items" 
	("house_model_id", "room_type", "amount")
VALUES 
	(1, 'ห้องนอน', '2'),
	(1, 'ห้องน้ำ', '2');
	
INSERT INTO "house_model_plans"
	("house_model_id", "floor", "size")
values
	(1, '1', '52 ตร.ม'),
	(1, '2', '58 ตร.ม');
	
INSERT into "house_model_plan_items"
	("house_model_plan_id", "room_type", "amount")
values
	(1, 'ห้องนอน', '1'),
	(1, 'ห้องน้ำ', '2'),
	(2, 'ห้องนอน', '2'),
	(2, 'ห้องน้ำ', '2');

INSERT INTO promotions
	("index", "heading", "description", "start_date", "end_date", "display")
VALUES 
	(1, 'Gareth Lang', 'No, you cant insert into multiple tables in one MySQL command. You can however use transactions.','2024/03/12', '2024/04/13', 'published');

INSERT INTO promotion_house_models 
	("promotion_id", "house_model_id")
VALUES
	(1, 1);
	
insert into promotion_free_items 
	("promotion_id", "description")
values
	(1, 'ฟรีโอน');

INSERT INTO interests (bank_name, interest_rate, note, display)
VALUES ('Bank A', '5', 'Interest rate for savings account', 'published'),
       ('Bank B', '4.5', 'Interest rate for fixed deposits', 'unpublished');

INSERT INTO activities ("index", heading, description, start_date, end_date, video_link, display)
VALUES (1, 'Activity 1', 'Description for Activity 1', '2024-05-01', '2024-05-31', 'https://youtube.com/activity1', 'published'),
       (2, 'Activity 2', 'Description for Activity 2', '2024-06-01', '2024-06-30', 'https://youtube.com/activity2', 'unpublished');
COMMIT;