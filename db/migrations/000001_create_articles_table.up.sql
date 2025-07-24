CREATE TABLE public.articles (
  id BIGSERIAL NOT NULL,
  title VARCHAR(100) NOT NULL,
  body TEXT NULL,
  status VARCHAR(20) NOT NULL,
  provider_type VARCHAR(50) NULL,
  link VARCHAR(255) NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ NULL,

  CONSTRAINT articles_pkey PRIMARY KEY (id),
  CONSTRAINT articles_status_check CHECK (status IN ('draft', 'published')),
  CONSTRAINT articles_provider_type_check CHECK (
    (provider_type IS NULL) OR
    (provider_type IN ('qiita', 'zenn', 'note'))
  )
) TABLESPACE pg_default;


CREATE INDEX IF NOT EXISTS idx_articles_status ON public.articles USING btree (status);
CREATE INDEX IF NOT EXISTS idx_articles_provider_type ON public.articles USING btree (provider_type);
CREATE INDEX IF NOT EXISTS idx_articles_created_at ON public.articles USING btree (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_articles_updated_at ON public.articles USING btree (updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_articles_provider_active ON public.articles (provider_type, status) WHERE deleted_at IS NULL;
