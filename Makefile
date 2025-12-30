.PHONY: buf-lint buf-fmt buf-gen buf

# buf lintチェック
buf-lint: _buf-exists
	buf lint --config buf.yaml

# buf フォーマット
buf-fmt: _buf-exists
	buf format -w --config buf.yaml

# buf コード生成
buf-gen: _buf-exists
	buf generate --template buf.gen.yaml --config buf.yaml

# buf関連タスクをまとめて実行
buf: buf-lint buf-fmt buf-gen

.SILENT:
_buf-exists:
ifeq ($(shell which buf),)
	@echo 'need buf. see https://docs.buf.build/installation'
	@exit 1
endif