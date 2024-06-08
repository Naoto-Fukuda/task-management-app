| # | Entity | UseCase | Parameters | Table/Index | API & Key Conditions |
|:-:|:-:|:-:|:-:|:-:|:-:|
|1|Tasks|createTask|{task}|Table|Query on PK|
|2|Tasks|getTaskById|{taskId}|Table|GetItem(PK = :taskId)|
|3|Tasks|getTasks||Table|Query on PK|
|4|Tasks|updateTaskById|{taskId}|Table|UpdateItem|
|5|Tasks|deleteTaskById|{taskId}|Table|DeleteItem|
|6|Tasks|getTasksByTitle|{title}|GSI-1|Query(GSI-1-PK  = :title)|
|7|Tasks|getTasksByDescription|{description}|GSI-1|Query(GSI-1-PK  = :description)|
|8|Tasks|getTasksByStatus|{status}|GSI-1|Query(GSI-1-PK  = :status)|
|9|Tasks|getTasksByTag|{tagName}|GSI-1|Query(GSI-1-PK  = :tagName)|
|10|Tasks|addTagToTask|{taskId, newTag}|Table|UpdateItem - Add newTag to Tag list|
|11|Tasks|updateTagOnTask|{taskId, oldTag, newTag}|Table|UpdateItem - Replace oldTag with newTag in Tag list|
|12|Tasks|deleteTagFromTask|{taskId, tagToDelete}|Table|UpdateItem - Remove tag from Tag list|
