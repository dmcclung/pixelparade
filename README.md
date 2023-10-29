modd

npx tailwindcss -i ./styles.css -o ../assets/styles.css --watch

docker compose run tailwind npx tailwindcss -c /src/tailwind.config.js -i /src/styles.css -o /dst/styles.css

docker build -t pixelparade .

docker run -it --rm --name pixelparade pixelparade

```bash
cd ~
mkdir -p apps/lenslocked.com
cd apps/lenslocked.com
git --work-tree=$HOME/apps/lenslocked.com \
    --git-dir=$HOME/repos/lenslocked.git \
    checkout -f main
ls
```