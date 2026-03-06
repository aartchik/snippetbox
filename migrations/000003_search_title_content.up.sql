ALTER TABLE snippets
  ADD FULLTEXT INDEX ft_snippets_title_content (title, content);