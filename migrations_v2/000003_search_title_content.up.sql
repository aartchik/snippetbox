CREATE INDEX ft_snippets_title_content
ON snippets
USING GIN (to_tsvector('simple', title || ' ' || content));
