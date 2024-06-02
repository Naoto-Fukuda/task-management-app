- シンプル

| PartionKey | Attributes |||||
|:-:|:-|:-|:-|:-|:-|
| TaskId | TaskName | Title| Description | Status|Tags|
| {TaskId} | {TaskName} | {Title} | {Description} | {Status} | {TagName1, ...} |
| ... | ... | ... | ... | ... | ... |

ただ、これでは懸念点がいくつか挙げれられます。
GetItemやQueryがEventIdにしか使えず、それ以外の属性で検索したい場合 Scan + Filterで高コストとなってしまう。
インデックスがない場合、項目数が少なくScanで問題ない場合が担保されている以外は、`Scan + Filter`は良くない検索方法として挙げられることが多い模様。
→ Scanではなく、QueryやGetItem / BatchGetItemの利用を検討するべき。

<br>

- シンプル + GSIパターン

| PrimaryKey | Attributes |||||
|:-:|:-|:-|:-|:-|:-|
| PartitionKey |GSI-1-PK|GSI-2-PK|GSI-3-PK|GSI-4-PK|GSI-5-PK|
| TaskId | TaskName | Title| Description | Status|Tags|
|||||
| {TaskId} | {TaskName} | {Title} | {Description} | {Status} | {TagName1, ...} |
| ... | ... | ... | ... | ... | ... |

各属性にGSIを取り入れることにより、Queryが可能になる。
しかし、これでも懸念点は残っており、GSIでは、別テーブルが作成されるためストレージコストが増加し、さらにスループットや管理コストが増加してしまう。

- GSI OverLoading

| PrimaryKey|PrimaryKey |Attributes|
|:-:|:-|:-|
| PK, GSI-1-SK |SK|GSI-1-PK|
| TaskId | DataType | DataValue | 
|||
| {TaskId} | TaskName | {TaskName} |
| {TaskId} | Title | {Title} |
| {TaskId} | Description |{Description} |
| {TaskId} | Status | {Status} |