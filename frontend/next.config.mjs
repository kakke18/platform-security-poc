import { readFile } from 'node:fs/promises';

// デプロイ時はパイプライン上で環境変数がprocess.envに設定される
// 開発時など、環境変数をローカル上で設定する場合はAPP_ENVを ./config/.env.xxx.json の拡張子に合わせて設定する
const env =
  process.env.NODE_ENV === 'production' && process.env.APP_ENV === 'prod'
    ? {}
    : {
        ...JSON.parse(
          // npm scriptsで設定したAPP_ENVに応じて読み込む.env.xxx.jsonを切り替える
          await readFile(`./config/.env.${process.env.APP_ENV}.json`),
        ),
      };

// process.env に設定してサーバー側（middleware、Server Components）でも使えるようにする
Object.assign(process.env, env);

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
};

export default nextConfig;
