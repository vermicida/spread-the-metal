# Spread The Metal

Spread The Metal es un bot que publica canciones de distintos subgéneros de rock y de metal en la cuenta de [Twitter](https://twitter.com/spreadthemetal) del mismo nombre. Lo hace de forma automática y con una cadencia de una hora: cada minuto 1 de cada hora del día se publica una nueva canción.

<img src="images/tweet.png" width="614" alt="Tweet">

## Cómo usar

Spread The Metal es una función [AWS Lambda](https://aws.amazon.com/lambda/) desarrollada en [Go](https://golang.org/). Debes tener instalado este lenguaje de programación en tu sistema para ejecutar los pasos que se proponen más abajo.

### Compilación y empaquetado

Lo primero que debemos hacer clonar este repositorio:

```bash
$ git clone https://github.com/vermicida/spread-the-metal.git
```

Con el código de la aplicación se distribuye un Makefile con varios targets que nos ayudarán con las siguientes tareas. Vamos a usar el primero de ellos para descargar las dependencias:

```bash
$ make deps
```

En este punto ya podemos compilar y empaquetar la aplicación. AWS Lambda espera el código de la función empaquetado en un documento zip, y es justo lo que hace el siguiente comando:

```bash
$ make dist
```

Al ejecutarlo, se compila el código y se comprime el binario resultante, dejando el resultado en `dist/function.zip`.

### Publicación

Como se ha comentado, Spread The Metal es una función Lambda. Se apoya en otros servicios de AWS, por lo que debemos crear en ellos los recusos oportunos para que todo funcione como se espera.

En primer lugar, vamos a crear una nueva tabla en [Amazon DynamoDB](https://aws.amazon.com/dynamodb/). El nombre de la tabla puede personalizarse, pues luego se configura en la función Lambda usando una variable de entorno, pero el nombre de las propiedades de los items sí que deben respetarse tal y como se define a continuación:

| Nombre | Key | Tipo | Descripción |
| :- | :- | :- | :- |
| Date | Partition Key | String | Almacena la fecha en que la canción debe publicarse |
| Hour | Sort Key | String | Almacena la hora en que la canción debe publicarse |
| Band | - | String | Nombre de la banda (se publica en el tweet) |
| Title | - | String | Título de la canción (se publica en el tweet) |
| Link | - | String | Link a la canción en Spotify (se publica en el tweet) |

Un ejemplo de item de esta tabla es el siguiente:

```json
{
    "Date": "20200315",
    "Hour": "19",
    "Band": "Sepultura",
    "Title": "Means to an End",
    "Link": "https://open.spotify.com/track/7IhrZIRIAlSMLJwO35dDQc?si=l5XiAAOFTgWa7yTb-_R5tA"
}
```

Atendiendo a las propiedades `Date` y `Hour` deducimos que el tweet de este item se publicará el 3  de marzo de 2020 a las 15:00 h.

La aplicación obtiene la hora actual con `time.Now()` para generar el Partition y Sort keys, y así obtener el item correspondiente a publicar. En caso de no existir registro para una fecha y hora concretas, no se publicaría nada en Twitter.

Lo siguiente que debe hacerse es habilitar en la función Lambda la posibilidad de obtener items de DynamoDB. Para esto debemos crear un nuevo rol y asociar una política que permita la acción `dynamodb:GetItem` sobre el ARN de la tabla de DynamoDB.

Y por fin, la función Lambda. Creamos una nueva, asignamos el rol recién creado en la sección **Execution role**, y establecemos lo siguiente en **Environment variables**:

| Variable | Requerida | Valor por defecto | Descripción |
| :- | :- | :- | :- |
| CONSUMER_KEY | Sí | | Consumer API key de nuestra aplicación de Twitter |
| CONSUMER_SECRET | Sí | | Consumer API secret key de nuestra aplicación de Twitter  |
| ACCESS_TOKEN | Sí | | Access token de nuestra aplicación de Twitter |
| ACCESS_SECRET | Sí | | Access token secret de nuestra aplicación de Twitter |
| STATUS_HASHTAGS | No | | Hashtags a incluir en el tweet (ej. #music #metal) |
| DEFAULT_REGION | Sí | eu-west-1 | Región donde está creada la tabla de canciones de DynamoDB |
| SONGS_TABLE_NAME | Sí | stm-songs | Nombre de la tabla de canciones de DynamoDB |
| DATE_FORMAT | Sí | 20060102 | Formato de fecha (debe corresponder con el elegido para el Partition Key de la tabla de canciones) |
| TIME_FORMAT | Sí | 15 | Formato de hora (debe corresponder con el elegido para el Sort Key de la tabla de canciones) |

Es importante indicar que debemos disponer de las credenciales oportunas de Twitter para que esta función Lambda pueda publicar tweets de forma autónoma. Para ello, creamos una nueva aplicación en nuestra cuenta [Developer](https://developer.twitter.com/en.html) de Twitter y actualizamos los valores de las variables de entorno afectadas.

Subimos el código de la función Lambda que hemos empaquetado previamente. Indicamos los siguientes valores en la sección **Function code**:

- **Code entry type:** Uplaod a .zip file
- **Runtime:** Go 1.x
- **Handler:** main
- **Function package:** seleccionamos el paquete zip

Guardamos los cambios y ya tendremos la función lista para empezar a operar.

Nos queda un último paso, que es automatizar la publicación de tweets cada hora; esto lo resolvemos con un [CloudWatch Event](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/WhatIsCloudWatchEvents.html). Lo hacemos en la sección **Designer**, añadiendo un nuevo trigger de tipo **CloudWatch Events/EventBridge** que tenga como regla la expresión `cron(1 * * * ? *)`.

Con esto ya estaría listo. Spread The Metal publicado y funcionando :-)
