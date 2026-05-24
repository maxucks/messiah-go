-- +goose Up
-- +goose StatementBegin
create table messages (
  id         uuid primary key default gen_random_uuid(),
  content    text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz,

  constraint check_content_not_empty
    check (trim(content) <> '')
);
  -- constraint fk_messages_user_id
  --   foreign key (user_id) references users (id) 
  --   on update cascade on delete set null
  -- constraint fk_messages_chan_id
  --   foreign key (chan_id) references channels (id)
  --   on update cascade on delete cascade
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table messages;
-- +goose StatementEnd
