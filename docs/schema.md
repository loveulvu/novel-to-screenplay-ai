# YAML Schema 设计

最终输出的剧本 YAML 以结构化、可校验、方便展示为目标。当前 MVP 先覆盖最小字段，后续可以继续扩展镜头、旁白、分场编号、道具和情绪节奏。

## 顶层字段

```yaml
title: "雨夜旧剧院"
source_chapters:
  - 1
characters:
  - id: "char_lin_che"
    name: "林澈"
    role: "主角"
scenes:
  - id: "scene_001"
    source_chapter: 1
    location: "老街雨巷"
    time: "夜晚"
    summary: "林澈收到父亲留下的铜钥匙。"
    characters:
      - "char_lin_che"
    dialogues:
      - character: "char_lin_che"
        emotion: "迟疑"
        line: "爸，如果这是你留下的线索，我一定会走到最后。"
    actions:
      - "林澈撑伞站在废弃邮筒前。"
```

## 字段定义

- `title`：剧本标题，必须非空。
- `source_chapters`：剧本覆盖的原小说章节编号，用于追溯来源。
- `characters`：全局角色列表，包含角色 `id`、名称和功能定位。
- `scenes`：剧本分场列表，是展示、复制、导出的核心内容。
- `scene.id`：分场唯一标识，便于校验、引用和后续局部重生成。
- `scene.source_chapter`：该分场主要来自的小说章节。
- `scene.location`：场景地点，必须非空。
- `scene.time`：场景时间，必须非空。
- `scene.summary`：分场摘要，必须非空。
- `scene.characters`：参与该场的角色 id 列表。
- `scene.dialogues`：台词列表。
- `scene.actions`：动作和舞台提示列表。
- `dialogue.character`：说话角色 id，必须非空。
- `dialogue.emotion`：台词情绪，用于帮助表演和后续镜头设计。
- `dialogue.line`：台词正文，必须非空。

## 为什么使用 title、source_chapters、characters、scenes

这四类字段组成最小可解释剧本结构。`title` 方便用户识别结果，`source_chapters` 让生成内容能回溯到原小说，`characters` 提供全局角色表，`scenes` 承载真正可展示和导出的剧本内容。

## 为什么角色使用 id

角色名可能重名、改名，也可能有别名。使用稳定的 `id` 可以避免后续生成、校验和局部编辑时混淆角色。展示层仍然可以把 `id` 映射回中文名。

## 为什么 dialogue 拆成 character、emotion、line

台词不只是文本。`character` 表明谁说，`emotion` 表明怎么说，`line` 表明说什么。拆开后可以支持校验、表演提示、镜头规划和后续更细的剧本格式导出。

## 为什么中间态用 JSON，最终输出 YAML

JSON 更适合作为后端接口和程序内部中间态，类型清晰，方便序列化、校验和前端消费。YAML 更适合作为最终交付结果，层级可读性更好，适合复制、编辑和人工讲解。
