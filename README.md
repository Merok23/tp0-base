### Resolución:
`make docker-compose up` para levantar la app.

Se agrego el volumen a el docker compose de forma que el directorio local ./client se corresponda a una carpeta config del cliente, en la cual se copia el archivo de config, y algo muy parecido en el servidor. De esta forma, si se cambia el archivo de configuración, no es necesario volver a construir las imágenes de docker, ya que el archivo de configuración se copia en el contenedor y se monta en el directorio correspondiente.

### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: `docker volumes`).