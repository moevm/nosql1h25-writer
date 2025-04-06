CREATE SCHEMA IF NOT EXISTS writer_sql;

CREATE  TABLE users ( 
	id                   serial  NOT NULL  ,
	email                varchar(256)  NOT NULL  ,
	public_name          varchar(64)  NOT NULL  ,
	"password"           varchar(256)  NOT NULL  ,
	balance              integer DEFAULT 0 NOT NULL  ,
	system_role          enum('user','admin') DEFAULT 'user' NOT NULL  ,
	is_active            boolean DEFAULT true NOT NULL  ,
	created_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	updated_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	CONSTRAINT pk_users PRIMARY KEY ( id ),
	CONSTRAINT unq_users UNIQUE ( email ) 
 );

CREATE  TABLE orders ( 
	id                   serial  NOT NULL  ,
	client_id            integer  NOT NULL  ,
	title                varchar(128)  NOT NULL  ,
	description          text  NOT NULL  ,
	completion_time      bigint  NOT NULL  ,
	cost                 integer    ,
	is_active            boolean DEFAULT true NOT NULL  ,
	freelancer_id        integer DEFAULT NULL   ,
	budget               integer DEFAULT 0 NOT NULL  ,
	created_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	updated_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	CONSTRAINT pk_orders PRIMARY KEY ( id )
 );

CREATE  TABLE profiles ( 
	id                   serial  NOT NULL  ,
	user_id              integer  NOT NULL  ,
	"role"               enum('client','freelancer')  NOT NULL  ,
	rating               real DEFAULT 0 NOT NULL  ,
	description          text  NOT NULL  ,
	created_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	updated_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	CONSTRAINT pk_profiles PRIMARY KEY ( id )
 );

CREATE  TABLE responses ( 
	order_id             integer  NOT NULL  ,
	freelancer_id        integer  NOT NULL  ,
	chat_id              integer  NOT NULL  ,
	is_active            boolean DEFAULT true NOT NULL  ,
	created_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	updated_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	CONSTRAINT pk_responses PRIMARY KEY ( order_id, freelancer_id, chat_id )
 );

CREATE  TABLE reviews ( 
	id                   serial  NOT NULL  ,
	author_id            integer  NOT NULL  ,
	profile_id           integer  NOT NULL  ,
	score                enum('1','2','3','4','5')  NOT NULL  ,
	content              text    ,
	created_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	CONSTRAINT pk_reviews PRIMARY KEY ( id )
 );

CREATE  TABLE statuses ( 
	order_id             integer  NOT NULL  ,
	sequential_number    integer  NOT NULL  ,
	title                enum('beginning','negotiation','budgeting','work','reviews','finished','dispute') DEFAULT 'beginning' NOT NULL  ,
	extra_info           jsonb    ,
	created_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	CONSTRAINT pk_statuses PRIMARY KEY ( order_id, sequential_number )
 );

CREATE  TABLE chats ( 
	id                   serial  NOT NULL  ,
	order_id             integer  NOT NULL  ,
	client_id            integer  NOT NULL  ,
	freelancer_id        integer  NOT NULL  ,
	is_active            boolean DEFAULT true NOT NULL  ,
	created_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	updated_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	CONSTRAINT pk_chats PRIMARY KEY ( id )
 );

CREATE  TABLE messages ( 
	id                   bigserial  NOT NULL  ,
	chat_id              integer  NOT NULL  ,
	sender_id            integer    ,
	is_system            boolean DEFAULT false NOT NULL  ,
	text                 text  NOT NULL  ,
	created_at           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL  ,
	CONSTRAINT pk_messages PRIMARY KEY ( id )
 );

CREATE INDEX idx_messages ON messages USING GIN ( text );

ALTER TABLE chats ADD CONSTRAINT fk_chats_users FOREIGN KEY ( client_id ) REFERENCES users( id );

ALTER TABLE chats ADD CONSTRAINT fk_chats_users_1 FOREIGN KEY ( freelancer_id ) REFERENCES users( id );

ALTER TABLE chats ADD CONSTRAINT fk_chats_orders FOREIGN KEY ( order_id ) REFERENCES orders( id );

ALTER TABLE messages ADD CONSTRAINT fk_messages_chats FOREIGN KEY ( chat_id ) REFERENCES chats( id );

ALTER TABLE messages ADD CONSTRAINT fk_messages_users FOREIGN KEY ( sender_id ) REFERENCES users( id );

ALTER TABLE orders ADD CONSTRAINT fk_orders_users FOREIGN KEY ( client_id ) REFERENCES users( id );

ALTER TABLE orders ADD CONSTRAINT fk_orders_users_1 FOREIGN KEY ( freelancer_id ) REFERENCES users( id );

ALTER TABLE profiles ADD CONSTRAINT fk_profiles_users FOREIGN KEY ( user_id ) REFERENCES users( id );

ALTER TABLE responses ADD CONSTRAINT fk_responses_users FOREIGN KEY ( freelancer_id ) REFERENCES users( id );

ALTER TABLE responses ADD CONSTRAINT fk_responses_orders FOREIGN KEY ( order_id ) REFERENCES orders( id );

ALTER TABLE reviews ADD CONSTRAINT fk_reviews_users FOREIGN KEY ( author_id ) REFERENCES users( id );

ALTER TABLE reviews ADD CONSTRAINT fk_reviews_profiles FOREIGN KEY ( profile_id ) REFERENCES profiles( id );

ALTER TABLE statuses ADD CONSTRAINT fk_statuses_orders FOREIGN KEY ( order_id ) REFERENCES orders( id );

COMMENT ON TABLE users IS 'Таблица пользователей.';

COMMENT ON COLUMN users.id IS 'ID пользователя.';

COMMENT ON COLUMN users.email IS 'Электронная почта пользователя, должна быть уникальной.';

COMMENT ON COLUMN users.public_name IS 'Публичное имя пользователя (ФИО или ФИ).';

COMMENT ON COLUMN users."password" IS 'Пароль в захэшированном виде.';

COMMENT ON COLUMN users.balance IS 'Количество средств на счету у пользователя.';

COMMENT ON COLUMN users.system_role IS 'Системная роль (пользователь, администратор, модератор и т.п.). Определяет возможность пользователя авторизоваться и попасть на служебные страницы.';

COMMENT ON COLUMN users.is_active IS 'Активен ли аккаунт пользователя в данный момент.';

COMMENT ON COLUMN users.created_at IS 'Дата создания записи.';

COMMENT ON COLUMN users.updated_at IS 'Дата обновления документа.';

COMMENT ON TABLE orders IS 'Заказы.';

COMMENT ON COLUMN orders.id IS 'ID заказа.';

COMMENT ON COLUMN orders.client_id IS 'ID пользователя, разместившего заказ.';

COMMENT ON COLUMN orders.title IS 'Название заказа.';

COMMENT ON COLUMN orders.description IS 'Описание заказа, которое задаётся пользователем.';

COMMENT ON COLUMN orders.completion_time IS 'Время, отведённое на выполнение заказа в наносекундах. Максимум при использовании int64 - 292 года.';

COMMENT ON COLUMN orders.cost IS 'Стоимость заказа. Если NULL - считать стоимость договорной. Используется только для показа в карточке заказа и никогда - для расчётов.';

COMMENT ON COLUMN orders.is_active IS 'Активен ли заказ (может быть неактивен, если заказ был скрыт администратором или удалён заказчиком).';

COMMENT ON COLUMN orders.freelancer_id IS 'ID выбранного исполнителя. Пока исполнитель не выбран, равен NULL.';

COMMENT ON COLUMN orders.budget IS 'Бюджет сделки, который резервируется в качестве гарантии оплаты услуг исполнителя. Может отличаться от cost.';

COMMENT ON COLUMN orders.created_at IS 'Дата создания записи.';

COMMENT ON COLUMN orders.updated_at IS 'Дата обновления записи.';

COMMENT ON TABLE profiles IS 'Профили пользователей.';

COMMENT ON COLUMN profiles.id IS 'ID профиля.';

COMMENT ON COLUMN profiles.user_id IS 'ID владельца профиля.';

COMMENT ON COLUMN profiles."role" IS 'Роль (заказчик или исполнитель).';

COMMENT ON COLUMN profiles.rating IS 'Рейтинг профиля, основан на отзывах.';

COMMENT ON COLUMN profiles.description IS 'Пользовательское поле "О себе".';

COMMENT ON COLUMN profiles.created_at IS 'Дата создания записи.';

COMMENT ON COLUMN profiles.updated_at IS 'Дата обновления записи.';

COMMENT ON TABLE responses IS 'Отклики исполнителей на заказы.';

COMMENT ON COLUMN responses.order_id IS 'ID заказа, на который оставлен отклик.';

COMMENT ON COLUMN responses.freelancer_id IS 'ID пользователя, оставившего отклик.';

COMMENT ON COLUMN responses.chat_id IS 'ID переписки (чата), созданного в рамках отклика.';

COMMENT ON COLUMN responses.is_active IS 'Активен ли отклик.';

COMMENT ON COLUMN responses.created_at IS 'Дата создания записи.';

COMMENT ON COLUMN responses.updated_at IS 'Дата редактирования записи.';

COMMENT ON TABLE reviews IS 'Отзывы пользователей.';

COMMENT ON COLUMN reviews.id IS 'ID отзыва.';

COMMENT ON COLUMN reviews.author_id IS 'ID пользователя, оставившего отзыв.';

COMMENT ON COLUMN reviews.profile_id IS 'ID профиля, на который оставлен отзыв.';

COMMENT ON COLUMN reviews.score IS 'Оценка (от 1 до 5).';

COMMENT ON COLUMN reviews.content IS 'Текст отзыва.';

COMMENT ON COLUMN reviews.created_at IS 'Дата создания записи.';

COMMENT ON TABLE statuses IS 'История изменения статусов заказов.';

COMMENT ON COLUMN statuses.order_id IS 'ID заказа, для которого присвоен статус.';

COMMENT ON COLUMN statuses.sequential_number IS 'Порядковый номер статуса по порядку.';

COMMENT ON COLUMN statuses.title IS 'Название статуса.';

COMMENT ON COLUMN statuses.extra_info IS 'Дополнительная информация о статусе, если нужно.';

COMMENT ON COLUMN statuses.created_at IS 'Дата создания записи.';

COMMENT ON TABLE chats IS 'Переписки.';

COMMENT ON COLUMN chats.id IS 'ID переписки.';

COMMENT ON COLUMN chats.order_id IS 'ID заказа, в рамках которого ведётся переписка.';

COMMENT ON COLUMN chats.client_id IS 'ID пользователя, выступающего в роли заказчика.';

COMMENT ON COLUMN chats.freelancer_id IS 'ID пользователя, выступающего в роли исполнителя.';

COMMENT ON COLUMN chats.is_active IS 'Является ли переписка активной (можно ли отправлять новые сообщения).';

COMMENT ON COLUMN chats.created_at IS 'Дата создания записи.';

COMMENT ON COLUMN chats.updated_at IS 'Дата обновления записи.';

COMMENT ON TABLE messages IS 'Сообщения в переписке.';

COMMENT ON COLUMN messages.id IS 'ID сообщения (при потоке в 1 млн. сообщений в секунду должно хватить на 292 тысячи лет).';

COMMENT ON COLUMN messages.chat_id IS 'ID переписки, к которой относится сообщение.';

COMMENT ON COLUMN messages.sender_id IS 'ID отправителя сообщения (может быть NULL), если сообщение системное.';

COMMENT ON COLUMN messages.is_system IS 'Является ли сообщение системным (системные сообщения должны отрисовываться по другому).';

COMMENT ON COLUMN messages.text IS 'Текст сообщения.';

COMMENT ON COLUMN messages.created_at IS 'Дата создания записи.';

