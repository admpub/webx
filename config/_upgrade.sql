ALTER TABLE official_film_item ADD FULLTEXT INDEX ft_actors_index (main_actors) WITH PARSER ngram;
ALTER TABLE official_film_item ADD FULLTEXT INDEX ft_title_index (title) WITH PARSER ngram;
ALTER TABLE official_film_item ADD FULLTEXT INDEX ft_directors_index (directors) WITH PARSER ngram;
ALTER TABLE official_film_item ADD FULLTEXT INDEX ft_screenwriters_index (screenwriters) WITH PARSER ngram;
ALTER TABLE official_film_item ADD FULLTEXT INDEX ft_areas_index (issued_areas) WITH PARSER ngram;
ALTER TABLE official_film_item ADD FULLTEXT INDEX ft_languages_index (languages) WITH PARSER ngram;
ALTER TABLE official_film_item ADD FULLTEXT INDEX ft_types_index (types) WITH PARSER ngram;
ALTER TABLE official_film_item ADD FULLTEXT INDEX ft_tags_index (tags) WITH PARSER ngram;