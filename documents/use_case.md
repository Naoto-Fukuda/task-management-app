|Entity|Usecase|Input Parameters|Return|Description|
|:----:|:-----:|:--------------:|:----:|:---------:|
|タスク|getTaskById|{"TaskId"}|{"TaskId", "Title", "Status", "Description", "Tags"}|task_idでタスクを一件取得|
|タスク|getTasks|{}|[{"TaskId", "Title", "Status", "Description", "Tags"}]|タスク一覧取得|
|タスク|getTasksByTag|{"Tag"}|[{"TaskId", "Title", "Status", "Description", "Tags"}]|タグでタスク一覧を取得|
|タスク|getTasksByTitle|{"Title"}|[{"TaskId", "Title", "Status", "Description", "Tags"}]|タイトルでタスク一覧を取得|
|タスク|getTasksByDescription|{"Description"}|[{"TaskId", "Title", "Status", "Description", "Tags"}]|説明でタスク一覧を取得|
|タスク|getTasksByStatus|{"Status"}|[{"TaskId", "Title", "Status", "Description", "Tags"}]|ステータスでタスク一覧を取得|
|タスク|createTask|{"Title", "Description", "Tags"}|{"TaskId", "Title", "Status", "Description", "Tags"}|タスクを作成|
|タスク|updateTaskById|{"TaskId", "Title", "Description", "Tags"}|{"TaskId", "Title", "Status", "Description", "Tags"}|task_idで指定されたタスクを更新|
|タスク|deleteTaskById|{"TaskId"}|{}|task_idで指定されたタスクを削除|
|タグ|addTagToTask|{"TaskId", "Tag"}|{"TaskId", "Title", "Status", "Description", "Tags"}|タスクにタグを追加|
|タグ|updateTagOnTask|{"TaskId", "OldTag", "NewTag"}|{"TaskId", "Title", "Status", "Description", "Tags"}|タスクに紐づくタグを更新|
|タグ|deleteTagFromTask|{"TaskId", "Tag"}|{"TaskId", "Title", "Status", "Description", "Tags"}|タスクからタグを削除|