- シンプル

| PartionKey | Attributes |||||
|:-:|:-|:-|:-|:-|:-|
| id | TaskName | Title| Description | Status|Tags|
| {TaskId} | {TaskName} | {Title} | {Description} | {Status} | {TagName1, ...} |
| ... | ... | ... | ... | ... | ... |

ただ、これでは懸念点がいくつか挙げられます。
GetItemやQueryがEventIdにしか使えず、それ以外の属性で検索したい場合 Scan + Filterで高コストとなってしまう。
インデックスがない場合、項目数が少なくScanで問題ない場合が担保されている以外は、`Scan + Filter`は良くない検索方法として挙げられることが多い模様。
→ Scanではなく、QueryやGetItem / BatchGetItemの利用を検討するべき。

<br>

- シンプル + GSIパターン

| PrimaryKey | Attributes |||||
|:-:|:-|:-|:-|:-|:-|
| PartitionKey |GSI-1-PK|GSI-2-PK|GSI-3-PK|GSI-4-PK|GSI-5-PK|
| id | TaskName | Title| Description | Status|Tags|
|||||
| {TaskId} | {TaskName} | {Title} | {Description} | {Status} | {TagName1, ...} |
| ... | ... | ... | ... | ... | ... |

各属性にGSIを取り入れることにより、Queryが可能になる。
しかし、これでも懸念点は残っており、GSIでは、別テーブルが作成されるためストレージコストが増加し、さらにスループットや管理コストが増加してしまう。

- GSI OverLoading

| PrimaryKey|Attributes|Attributes|
|:-:|:-|:-|
| PK, GSI-1-SK |SK|GSI-1-PK|
| ID | DataType | DataValue | 
|||
| {TaskId} | TaskName | {TaskName} |
| {TaskId} | Title | {Title} |
| {TaskId} | Description |{Description} |
| {TaskId} | Status | {Status} |
| {TaskId} | Tag | {TagName1, TagName2, ...} |

これでインデックスを使用し、scanせずに検索が可能となるパターンは以下の通りです。
1.**TaskId(PK)での検索**:タスク一覧表示
```
主キー（`PK`）が`TaskId`であるため、`GetItem`や`Query`操作を使用して、特定の`TaskId`を持つタスクのすべてのデータを取得します。
```
例えば、`TaskId=1`を指定して、取得することで、タスク名、タイトル、説明、ステータス、タグなど、そのタスクIDに関連付けられたすべての属性を取得できる。

2.**TaskIdとDataTypeで(SK)の検索**:個別アイテムのタイトルやステータスの更新
```
`PK`が`TaskId`で、`SK`が`DataType`であるため、`Query`操作を使用して、特定の`TaskId`と`DataType`の組み合わせを持つデータを取得します。
```
例えば、`TaskId=1`とDataType=Titleを指定して、取得することで、そのタスクのタイトルのみを取得できる。

3.**DataValue(GSI-1-PK)とDataType(SK)での検索**:タグやタイトルで検索
```
`GSI-1-PK`が`DataValue`で、`SK`が`DataType`であるため、`Query`操作を使用して、特定の`DataValue`と`DataType`の組み合わせを持つデータを取得します。
```
例えば、DataTypeがStatusでDataValueが`complete`の場合、ステータスが"complete"であるすべてのタスクの情報を取得できる。

今回の要件はおそらく上記3つで満たせそうですが、副次的に以下の検索をインデックスを使用して可能となります。
4.**DataValue(GSI-1-PK)での検索**:
```
`GSI-1-PK`が`DataValue`であるため、`Query`操作を使用して、特定の`DataValue`を持つデータを取得します。
```
例えば、`DataValue=complete`を指定して、取得することで、`complete`ステータスを持つアイテムを取得することができる。ただし、3で検索するほうが効率的に検索を行うことが可能なので使用しない。

5.**DataValue(GSI-1-PK)とTaskID(GSI-1-SK)での検索**:
```
`GSI-1-PK`が`DataValue`で、GSI-1-SKが`TaskId`であるため`Query`操作を使用して、特定の`DataValue`と`TaskId`の組み合わせを持つデータを取得します。
```
例えば、DataValueが特定のタグ名で、TaskIdが特定のタスクIDであれば、そのタグを持つ特定のタスクの情報を取得できる。
