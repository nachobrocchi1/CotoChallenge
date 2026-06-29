# Challenge para Coto.

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
Tambien al tomar la decision de usar un slice en memoria me obliga a asignarle un Id a cada Sale y validar que no exista previamente.

- Se resolvio el calculo del tiempo de ejecucion de cada operacion usando un middleware en los handler evitando código duplicado.

- Se crearon funciones genericas en la capa de transporte http para hacer Encode y Decode de requests y response evitando codigo duplicado.

- Creacion de tests con asistencia de AI para optimizar tiempo.

- Comandos complejos del Makefile optimizados con AI.

- Mocks y Repository fueron excluidos de la cobertura de los tests ya que no aporta valor.
