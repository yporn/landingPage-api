BEGIN;

INSERT INTO "roles" ("title") VALUES ('dev'), ('admin');

INSERT INTO
    "users" (
        "username", "name", "email", "password", "role_id"
    )
VALUES (
        'y.pornwisa', 'pornwisa', 'y.pornwisa@gmail.com', '$2a$10$8KzaNdKIMyOkASCH4QvSKuEMIY7Jc3vcHDuSJvXLii1rvBNgz60a6', 2
    ),
    (
        'admin', 'admin', 'admin001@kawaii.com', '$2a$10$3qqNPE.TJpNGYCohjTgw9.v1z0ckovx95AmiEtUXcixGAgfW7.wCi', 2
    );

INSERT INTO
    "careers" (
        "position", "amount", "location", "description", "qualification", "start_date", "end_date", "status", "display"
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

COMMIT;