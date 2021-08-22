# API Scrapping de Sismología

## Sismos Recientes

Permite obtener una lista de los sismos recientes:

Ejemplo: https:/api.xor.cl/sismo/recent

Es posible filtrar por magnitud con el parámetro GET `magnitude`, entregando sismos de esa magnitud o mayor

Ejemplo: https://api.xor.cl/sismo/recent?magnitude=5
## Sismos históricos por día

Permite obtener una lista de los sismos ocurridos un día (en horario local) usando la versión YYYYMMDD de la fecha a buscar.

Ejemplo: https:/api.xor.cl/sismo/historic/20100227

Es posible filtrar por magnitud con el parámetro GET `magnitude`, entregando sismos de esa magnitud o mayor

Ejemplo: https://api.xor.cl/sismo/historic/20100227?magnitude=5
