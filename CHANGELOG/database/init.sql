create table annotation_templates (
    id varchar(64) not null,
    name varchar(128) not null,
    type varchar(64) not null,
    create_at bigint not null,
    -- soft delete
    delete_at bigint default 0,
    update_at bigint not null,
    description varchar(128),
    PRIMARY KEY (id)
);

create table annotation_template_exts (
    annotation_template_id varchar(64) not null,
    word_list text,
     -- soft delete
    delete_at bigint default 0,
    CONSTRAINT anno_template_ukey UNIQUE (annotation_template_id, delete_at)
);

create table labels (
    id varchar(64) not null,
    annotation_template_id varchar(64) not null,
    name varchar(64) not null,
    super_category_name varchar(64),
    color int default 0,
    create_at bigint not null,
    -- soft delete
    delete_at bigint default 0,
    update_at bigint not null,
    key_point_def varchar(512),
    key_point_skeleton varchar(512),
    cover_image_url varchar(512),
    PRIMARY KEY (id),
    CONSTRAINT label_ukey UNIQUE (annotation_template_id, name, delete_at)
);

create table data_views (
    id varchar(64) not null,
    related_data_view_id varchar(64),
    annotation_template_id varchar(64),
    name varchar(64) not null,
    view_type varchar(64) not null,
    raw_data_type varchar(64) not null,
    progress float default 0,
    status varchar(64),
    commit_id varchar(64),
    zip_format varchar(64),
    raw_data_view_id varchar(64),
    annotation_view_id varchar(64),
    train_raw_data_view_id varchar(64),
    train_annotation_view_id varchar(64),
    val_raw_data_view_id varchar(64),
	val_annotation_view_id varchar(64),
    description varchar(128) default "",
    create_at bigint not null,
    delete_at bigint default 0,
    -- soft delete
    update_at bigint not null,
    PRIMARY KEY (id)
);

create table data_view_items (
    data_view_id varchar(64) not null,
    data_item_id varchar(64) not null,
    PRIMARY KEY (data_view_id, data_item_id)
);

create table annotations (
    id varchar(64) not null,
    data_item_id varchar(64) not null,
    annotation_template_id varchar(64) not null,
    text_data varchar(1024),
    PRIMARY KEY (id)
);
create index anno_data_item_idx on annotations (data_item_id)
create index anno_anno_temp_idx on annotations (annotation_template_id)

create table data_items (
    id varchar(64) not null,
    name varchar(128) not null,
    type varchar(32) not null,    
    create_at bigint not null,
    PRIMARY KEY (id)
);

create table image_exts (
    id varchar(64) not null,
    thumbnail varchar(64),
    sha256 varchar(128) not null,
    size bigint,
    width int,
    height int,
    PRIMARY KEY (id),
    CONSTRAINT image_sha256 UNIQUE (sha256)
);

create table image_scores (
    id varchar(64) not null,
    light float,
    dense float,
    shelter float,
    size float,
    PRIMARY KEY (id)
);

create table rgbd_exts (
    id varchar(64) not null,
    sha256 varchar(128) not null,
    image_size bigint,
    image_width int,
    image_height int,
    depth_size bigint,
    depth_width int,
    depth_height int,
    PRIMARY KEY (id),
    CONSTRAINT rgbd_sha256 UNIQUE (sha256)
);

create table points_3d_exts (
    id varchar(64) not null,
    sha256 varchar(128) not null,
    size bigint,
    xmin float,
    xmax float,
    ymin float,
    ymax float,
    zmin float,
    zmax float,
    rmean float,
    gmean float,
    bmean float,
    rstd float,
    gstd float,
    bstd float,
    PRIMARY KEY (id),
    CONSTRAINT points3d_sha256 UNIQUE (sha256)
)

create table model_exts (
    id varchar(64) not null,
    sha256 varchar(128) not null,
    size bigint,
    PRIMARY KEY (id)
);

create table locations (
    data_item_id varchar(64) not null,
    bucket_name varchar(128),
    object_key varchar(512),
    name varchar(64),
    environment varchar(64),
    CONSTRAINT loc_ukey UNIQUE (data_item_id, name, environment)
);

create table raw_data_labels(
    raw_data_id varchar(64) not null,
    annotation_id varchar(64) not null,
    label_id varchar(64) not null
);
create index raw_data_labels_anno_id_idx on raw_data_labels (annotation_id)
create index raw_data_labels_label_id_idx on raw_data_labels (label_id)

-- for shed lock
CREATE TABLE shedlocks (
  name VARCHAR(64),
  lock_until bigint not NULL,
  locked_at bigint not  NULL,
  locked_by VARCHAR(255),
  PRIMARY KEY (name)
);

