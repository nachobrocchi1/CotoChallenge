# Challenge Coto.

## Como ejecutar
### Requisitos
- Go 1.26

### Ejecutar Local
```make run``` ejecuta localmente el servidor en el puerto 8080. 
En el archivo Makefile se puede modificar tambien el puerto en el que se quiere iniciar el servidor actualizando la variable PORT.
Tambien se puede ejecutar directamente usando PORT=8080 go run cmd/api/main.go desde el directorio raiz.

#### POSTMAN
Importar el archivo [CotoChallenge.postman_collection.json](CotoChallenge.postman_collection.json) a Postman para probar la aplicacion

#### CURL
- Insertar Venta
``` bash 
curl --location 'localhost:8080/api/v1/sales' \
--header 'Content-Type: application/json' \
--data '{
    "vehicle": "sedan",
    "center": "center1"
}'
```

- Obtener el volumen de ventas total.
``` bash 
curl --location 'localhost:8080/api/v1/volume'
```

- Obtener el volumen de ventas por centro.
``` bash 
curl --location 'localhost:8080/api/v1/volume/center'
```

- Obtener el volumen de ventas por centro.
``` bash 
curl --location 'localhost:8080/api/v1/sales/percentage'
```

### Ejecutar tests
```make test```

### Ver cobertura
```make coverage```


## Estructura
Se utilizo una arquitectura Clean Architecture, organizando el codigo en capas donde las dependencias apuntan siempre hacia adentro (las capas externas solo conocen a las internas, pero no al reves).

```internal/domain``` contiene las definiciones de entidades de negocio. Tambien incluye reglas de negocio como por ejemplo el impuesto agregado a los vehiculos de tipo Sport y el calculo del precio final de cada vehiculo.

```internal/repository``` Define abstraccion de persistencia mediante interface ```SaleRespository```. Incluye una implementacion concreta ```in_memory_repository.go``` que contiene un slice protegido por un ```sync.RWmutex()``` para garantizar accesos seguros entre requests.

```internal/service``` Define los casos de uso del sistema. Es donde vive la logica de negocio.

```internal/transport/http``` Capa mas externa. Responsable de decodificar una peticion, validar campos de entrada para llamar a la capa service. Serializa la respuesta.
El subdirectorio ```midleware/``` implementa tareas genericas entre los endpoints.

```cmd/api``` Punto de entrada de la aplicacion. Solo 'conecta' las capas instanciando dependencias en orden e inicia el servidor.

## Decisiones tecnicas

- Para el servidor se eligio la libreria estandar ```net/http``` ya que se puede crear endpoints ```METHOD /path```de forma nativa usando ```ServeMux``` reduciendo complejidad y dependencias.

- Persistencia en memoria. Para no añadir complejidad a la solucion se eligio persistir la informacion en memoria utilizando un Slice y protegiendo accesos usando ```sync.RWMutex```.
Este tipo de Mutex permite hacer ```RLock()``` para hacer lecturas simultaneas (multiples requests de lectura no se bloquean entre sí) y ```Lock()``` para escritura donde nadie mas puede escribir o leer el slice.
Se utilizó el paquete uuid de google para generar el Id, es practicamente imposible que se genere un Id repetido por lo cual no se añadió validacion de Id al agregar.
Al utilizar un ```interface``` para el repository, se puede realizar distintas implementaciones del mismo y añadir por ejemplo una base de datos PostgreSQL sin tener que modificar ninguna otra capa de la solucion.

- Se resolvio el calculo del tiempo de ejecucion de cada operacion usando un middleware en los handler evitando código duplicado. Tambien se incluyó un middleware 'Recovery' para recuperacion después de un ```panic```, para no dejar un request sin respuesta en caso de un error grave.

- Se crearon funciones genericas en la capa de transporte http para hacer Encode y Decode de requests y response evitando codigo duplicado.

- Creacion de tests con asistencia de AI para optimizar tiempo.

- Comandos complejos del Makefile optimizados con AI.

- Los paquetes ```mocks```, ```repository``` y ```main``` fueron excluidos de la cobertura de los tests ya que no tienen logica testeable.

- Se crearon structs diferenciados por capa, por ejemplo CreateSaleRequest vs domain.Sale. El usuario no deberia poder enviar Id, Price o Date.

- Validacion de request con mensaje explicito indicando todos los errores en un solo mensaje en lugar de fallar a medida que el usuario corrige su request.

- Se añadio verificacion de linter para mantener la calidad del código.

- No se implementó paginación ya que los endpoints de volumen devuelven datos agregados acotados por la cantidad de centros de distribución, lo cual no justifica su uso.

- Trabajando en el endpoint /api/v1/sales/percentage se decidió crear un paquete DTO (data transfer object) que comparte structs tanto en la capa de transporte como en la capa service. Esto surge por la necesidad de utilizar structs en ambas capas ya sea para las operaciones en la capa Service, como para el struct de respuesta del endpoint.
Siendo muy estricto yo hubiese duplicado los structs en ambas capas, o por el contrario, crearlo en una sola capa y utilizarlo en ambos lugares, en este caso se eligió una tercera opción, donde los structs viven en un paquete DTO y pueden ser importados por ambas capas.

## Estimacion

### Tareas
- Setup del proyecto
- Flujo de insertar una venta
- Flujo de obtener ventas totales
- Flujo de obtener ventas por centro.
- Flujo de obtener el porcentaje de unidades de cada modelo vendido en cada centro sobre el total de ventas.

### Esfuerzo
Las tareas mas pesadas son el setup del proyecto, ya que es crear la estructura y tomar las decisiones de arquitectura, y el flujo de porcentage de unidades en cada centro sobre total de ventas, ya que podría ser complejo desarrollarlo, el resto son tareas simples una vez que el setup del proyecto esta terminado. Yo estimaria 1 punto a esas dos tareas principales, y luego 1 punto entre flujos de insertar, obtener volumen total y obtener volumen por centro. Suma total de 3 puntos, lo cual yo mediría en 1 o 2 dias de trabajo o entre 8 y 16 horas.
