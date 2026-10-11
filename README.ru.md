<div align="center">

<picture>
  <source media="(prefers-color-scheme: dark)" srcset=".github/assets/banner-ru-dark.svg">
  <img alt="Mikan 2.0 — открытый форк VPN-панели на mihomo" src=".github/assets/banner-ru-light.svg" width="100%">
</picture>

<br>

**Mikan 2.0 — открытый форк Mikan с бесплатной дорожной картой и собственным релизным процессом.**

[![Репозиторий](https://img.shields.io/badge/GitHub-Mikan%202.0-ff853d?style=flat-square&logo=github)](https://github.com/m0ruy0/mikan-2.0)
[![Лицензия](https://img.shields.io/badge/license-GPL--3.0-2f9e6b?style=flat-square)](LICENSE)
[![Ядро](https://img.shields.io/badge/core-mihomo-9e8eff?style=flat-square)](https://github.com/MetaCubeX/mihomo)

**Русский** · [English](README.md)

</div>

---

## Установка одной командой

> **Первый релиз Mikan 2.0 ещё не опубликован.** Команда установки приведена для будущего стабильного релиза; пока не запускайте её на production-сервере.

На чистом VPS с Ubuntu 22.04+ или Debian 12+ (amd64/arm64), после публикации первого подписанного релиза:

```bash
curl -fsSL https://github.com/m0ruy0/mikan-2.0/releases/latest/download/install.sh | sudo bash
```

Установщик будет проверять подпись манифеста и SHA-256 загруженного бинарника, затем запускать мастер настройки. Не используйте для Mikan 2.0 установщик или Docker-образ оригинального Mikan.

Подробности: [инструкция по установке](https://m0ruy0.github.io/mikan-2.0/start/install/).

## Что развиваем

- **Ноды Pterodactyl** — Egg для обычного игрового хостинга без root и Docker socket внутри контейнера. Совместимость будет зависеть от TCP/UDP allocations провайдера.
- **Магазин аддонов** — несколько готовых Telegram-ботов, платёжные адаптеры и интеграции с разрешениями и проверяемыми обновлениями.
- **API для собственных ботов** — документированный API, ключи с ограниченными правами и подписанные webhooks.
- **Бесплатные возможности** — перенос панели, безопасное подключение внешних подписок как нод и отдельный админ-бот Telegram.

Все эти возможности пока относятся к дорожной карте, а не к готовому релизу. Полный план: [docs/mikan20-roadmap.md](docs/mikan20-roadmap.md).

## Разработка

Основа унаследована от [Mikan](https://github.com/Miroshka000/mikan); исходные уведомления об авторстве сохранены. Mikan 2.0 — модифицированный форк под GNU GPL v3, не официальный релиз авторов исходного проекта.

Разработка идёт на ветке `main`. При изменении установщика необходимо проверять Go-тесты, сборку web/site и Rust-тесты; релиз не считается готовым, пока все проверки не пройдены.

## Лицензия

[GNU GPL v3](LICENSE). Коммерческое использование и распространение возможны на условиях лицензии.
