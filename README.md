# hacku2020
お祈りメールが来たときに励ましてくれるシステム．

## 前準備
### 環境変数
- `EMAIL_SERVER`
    - メールサーバのIPアドレスおよびポート番号
- `EMAIL_ADDR`
    - 監視対象のメールアドレス
- `EMAIL_PASSWORD`
    - メールボックスにアクセスするためのパスワード
- `EMAIL_BOX`
    - 監視対象のメールボックス 
- `SLACK_TOKEN`
    - Slackボットのトークン
- `SLACK_CHANNEL`
    - Slackの投稿先チャンネルのID

### 応援キャラの設定
例)
```yaml
characters:
- name: "はげましちゃん1"
  icon: ":hagemashi1:"
  encourage-message:
  - "これからです！"
  - "今回はたまたまあなたの実力が伝わらなかっただけですよ"
  - "まだまだ落ち込むところじゃありませんよ！"
  - "よく頑張りましたね"
  - "今日はゆっくり休んで明日から頑張りましょう！"
  - "私はいつでもあなたを応援してますよ"
  congratulatory-message:
  - "おめでとうございます！"
- name: "はげましちゃん2"
  icon: ":hagemashi2:"
  encourage-message:
  - "まあそんなときもあるでしょ 次回頑張んなよ"
  - "あんたがいつも頑張ってんのあたしは知ってるよ"
  congratulatory-message:
  - "おめでとう　よかったじゃん"
```

- `name`
    - キャラクターの名前
- `icon`
    - キャラクターのアイコン（slackに絵文字として登録しているもののみ使用可能）
- `encourage-message`
    - お祈りメールが届いた時に応援してくれるメッセージのリスト
- `congratulatory-message`
    - 合格したときのお祝いメッセージ
    
## 使い方
```bash
# build application
$ go build -o ./oinori

# run application
$ ./oinori --character-config-path ./character_config.yaml \
           --message-num 4 \
           --light-host 127.0.0.1 \
           --m5stack-host m5stack.local
```

### オプション
- `--character-config-path`
    - キャラクターの設定ファイルへのパス
- `--message-num`
    - 一度の励ましorお祝いで投稿するメッセージ数
    - 用意されているメッセージより多い数が指定されている場合，用意されている数に合わされる
- `--light-host`
    - LightAPIを提供しているホストのIPまたはドメイン名
- `--m5stack-host`
    - M5stackのIPアドレスまたはドメイン名

## 関連リポジトリ
- [hacku2020_m5stack](https://github.com/StudioAquatan/hacku2020_m5stack)
- [HackU2020_AR](https://github.com/StudioAquatan/HackU2020_AR)
