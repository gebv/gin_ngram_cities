# Пример ngram+GIN для полнотекстового поиска

Подготовка данных
1. Скачать https://github.com/x88/i18nGeoNamesDB/blob/master/i18n_GeoCSVDump_with_header_qouted_delimiter_comma_v0.4.zip
2. Только английский `cat _cities.csv | cut -d ',' -f14 > _cities_en.csv`


```
2022/07/10 15:49:34 total cities: 2246419
2022/07/10 15:49:38 total records in GIN: 14782
2022/07/10 15:49:38
2022/07/10 15:49:38 lookup: Saint-Petersburg
2022/07/10 15:49:38 matched count: 4
2022/07/10 15:49:38 result: []string{"Saint Petersburg", "Saint Petersburg", "Borough of Saint Petersburg", "Saint Petersburg"}
2022/07/10 15:49:38
2022/07/10 15:49:38 lookup: St Petersburg
2022/07/10 15:49:38 matched count: 4
2022/07/10 15:49:38 result: []string{"Saint Petersburg", "Saint Petersburg", "Borough of Saint Petersburg", "Saint Petersburg"}
```

NOTE: в текущей реализации используются только [a-zA-Z0-9] символы

TODO:
- ранжирование результата поиск
- улучшить санитацзинг слов
