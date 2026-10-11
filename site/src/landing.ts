// Texts of the landing page and the footer in both languages. Docs pages live in
// src/content/docs (Russian) and src/content/docs/en (English).

export const links = {
	github: 'https://github.com/m0ruy0/mikan-2.0',
	releases: 'https://github.com/m0ruy0/mikan-2.0/releases',
	issues: 'https://github.com/m0ruy0/mikan-2.0',
	license: 'https://github.com/m0ruy0/mikan-2.0/blob/main/LICENSE',
	sponsor: 'https://m0ruy0.github.io/mikan-2.0/support/',
	mihomo: 'https://github.com/MetaCubeX/mihomo',
};

// Two lines, so the command fits the hero without wrapping in odd places; bash reads the
// trailing backslash as one command when it is pasted.
export const installCommand =
	'curl -fsSL https://github.com/m0ruy0/mikan-2.0/releases/latest/download/install.sh \\\n  | sudo bash';

export type Lang = 'ru' | 'en';
export type Screen = 'dashboard' | 'users' | 'inbounds' | 'telegram' | 'settings';

const docs = (lang: Lang, path: string) => `${import.meta.env.BASE_URL}${lang === 'en' ? 'en/' : ''}${path}`;

export const strings = {
	ru: {
		hero: {
			title: 'Mikan 2.0 — VPN-панель с открытым развитием',
			lead: 'Открытый форк Mikan: сохраняем основу панели и развиваем ноды Pterodactyl, внешний источник подписок, магазин аддонов и API для собственных ботов. Новые возможности планируются бесплатными.',
			requirements: 'Поддерживаемая установка появится с первым подписанным релизом. До публикации релиза не запускайте команду на production-сервере.',
			primary: 'Установка',
			secondary: 'Дорожная карта',
			art: 'Новый символ Mikan 2.0 — мандарин с сетевыми орбитами',
		},
		screens: {
			title: 'Одна панель на всё',
			lead: 'Реальные экраны исходного интерфейса Mikan. По мере развития форка снимки будут обновляться.',
			label: 'Экраны панели',
			items: {
				dashboard: ['Обзор', 'Онлайн, трафик за сутки и состояние каждого подключения на одном экране.'],
				users: ['Пользователи', 'Тарифы, сроки, трафик и устройства каждого подписчика, фильтры и действия над группой.'],
				inbounds: ['Подключения', 'Протоколы с портами и маскировкой. Видно, какие приложения получают каждый из них.'],
				telegram: ['Telegram', 'Бот настраивается прямо в панели: меню, тексты, уведомления и рассылки, с предпросмотром.'],
				settings: ['Настройки', 'Подписка, маршрутизация для Clash, безопасность входа и импорт из других панелей.'],
			} satisfies Record<Screen, [string, string]>,
		},
		features: {
			title: 'Что уже есть в основе',
			items: [
				{
					title: 'Ядро mihomo внутри ноды',
					text: 'mihomo встроен как библиотека, поэтому трафик считается по каждому соединению каждого пользователя. Тарифы по сроку и гигабайтам, лимит устройств и привязка к ним.',
					href: docs('ru', 'start/overview/'),
				},
				{
					title: 'Пятнадцать протоколов',
					text: 'VLESS REALITY с Vision, XHTTP и gRPC, Hysteria2, TUIC, AnyTLS, Trojan, TrustTunnel, Mieru и другие. Каждый добавляется из пресета, ключи создаются сами.',
					href: docs('ru', 'protocols/presets/'),
				},
				{
					title: 'Сама переживает блокировки',
					text: 'Если клиенты из разных сетей перестали доходить до протокола, панель переносит его на свободный HTTPS-порт или подбирает новый сайт маскировки рядом с сервером.',
					href: docs('ru', 'protocols/autotune/'),
				},
				{
					title: 'Ноды, каскады и WARP',
					text: 'Удалённая нода ставится одной командой с ключом и обновляется из панели. Протокол может выходить в интернет через другую ноду или через Cloudflare WARP.',
					href: docs('ru', 'nodes/nodes/'),
				},
				{
					title: 'Подписка под каждое приложение',
					text: 'Подписка узнаёт Happ, v2RayTun, Koala Clash, Clash Verge, Hiddify и другие и отдаёт каждому только то, что оно запускает. Для Clash есть маршруты по сервисам.',
					href: docs('ru', 'subscriptions/page/'),
				},
				{
					title: 'Telegram-бот и Mini App',
					text: 'Подписчики смотрят срок и устройства, покупают и продлевают доступ. Оплата через Telegram Stars и адаптеры маркетплейса: ЮKassa, CryptoBot.',
					href: docs('ru', 'billing/telegram/'),
				},
				{
					title: 'Промокоды, пулы и пакеты',
					text: 'Скидки, бонусные дни и трафик. Отдельный лимит для части протоколов, докупка гигабайтов в боте, пробный период.',
					href: docs('ru', 'billing/promo-codes/'),
				},
				{
					title: 'Переезд с других панелей',
					text: 'Импорт пользователей из Marzban, PasarGuard и Remnawave. Старые ссылки подписок работают, когда старый домен указывает на mikan.',
					href: docs('ru', 'operations/import/'),
				},
				{
					title: 'Спокойные обновления',
					text: 'PostgreSQL, подписанные релизы, бэкап перед каждым обновлением и откат, если новая версия не поднялась. Бэкап базы может каждый день приходить в Telegram.',
					href: docs('ru', 'operations/updates/'),
				},
			],
		},
		protocols: {
			title: 'Каждому приложению то, что оно запускает',
			lead: 'Подписка определяет приложение и не отдаёт ему протоколы, которые оно не поднимет. Никаких «у меня не работает» из-за неподходящего ключа.',
			headers: ['Протокол', 'Clash-приложения', 'Happ, v2RayTun', 'Hiddify, Karing'],
			yes: 'есть',
			no: 'нет',
			note: '¹ Один ключ на всех: учёта и лимитов по пользователям на этих протоколах нет.',
			more: 'Все пресеты и их настройки',
			moreHref: docs('ru', 'protocols/presets/'),
		},
		phones: {
			title: 'Страница подписки и Mini App',
			text: 'Подписчик открывает ссылку в браузере или в Telegram и видит срок, трафик и устройства. Платформа определяется сама, а кнопка «Добавить» открывает подписку сразу в приложении: Happ, v2RayTun, Koala Clash, Hiddify и другие.',
			link: 'Про подписки и приложения',
			href: docs('ru', 'subscriptions/page/'),
			alts: ['Страница подписки на телефоне', 'Обзор панели на телефоне'],
		},
		roadmap: {
			badge: 'Бесплатная дорожная карта',
			title: 'Что развиваем в Mikan 2.0',
			lead: 'Эти функции планируются для бесплатной основной версии. Начало разработки не зависит от сбора пожертвований.',
			all: 'Полная дорожная карта',
		},
		cta: {
			title: 'Mikan 2.0 — открытая разработка',
			text: 'Стабильная установка одной командой будет доступна после первого подписанного релиза. Пока можно изучить код и дорожную карту.',
			docs: 'Инструкция по установке',
			sponsorTitle: 'Разработка без закрытых функций',
			sponsorText: 'Базовые возможности из дорожной карты планируются бесплатными. Статус каждой идеи и технические ограничения опубликованы открыто.',
			sponsor: 'Открыть дорожную карту',
			goals: 'Планы разработки',
			github: 'Исходный код на GitHub',
		},
		header: {
			label: 'Проект',
			githubLabel: 'GitHub: исходный код mikan',
			support: 'Дорожная карта',
			supportLabel: 'Дорожная карта Mikan 2.0',
			otherLanguage: 'English version',
		},
		footer: {
			label: 'Ссылки проекта',
			license: 'Mikan 2.0 — модифицированный форк под GNU GPL v3; сохраняется авторство исходного проекта.',
			releases: 'Релизы',
			issues: 'Сообщить о проблеме',
			sponsor: 'Дорожная карта',
		},
	},
	en: {
		hero: {
			title: 'Mikan 2.0 — a VPN panel built in the open',
			lead: 'An open Mikan fork: we keep the panel foundation and develop Pterodactyl nodes, external subscription sources, an addon store and an API for custom bots. New core features are planned to remain free.',
			requirements: 'The supported installer will be available with the first signed release. Do not run the command on a production server before that release is published.',
			primary: 'Installation',
			secondary: 'Roadmap',
			art: 'Mikan 2.0 mandarin mascot with network orbits',
		},
		screens: {
			title: 'One panel for all of it',
			lead: 'Real screens from the upstream Mikan interface. Screenshots will be refreshed as this fork evolves.',
			label: 'Panel screens',
			items: {
				dashboard: ['Overview', 'Who is online, traffic for the day and the state of every protocol on one screen.'],
				users: ['Users', 'Plans, terms, traffic and devices of every subscriber, with filters and bulk actions.'],
				inbounds: ['Protocols', 'Protocols with their ports and camouflage, and which apps get each of them.'],
				telegram: ['Telegram', 'The bot is set up right in the panel: menu, texts, notices and broadcasts, with a preview.'],
				settings: ['Settings', 'Subscription, routing for Clash apps, sign-in security and imports from other panels.'],
			} satisfies Record<Screen, [string, string]>,
		},
		features: {
			title: 'What the base already includes',
			items: [
				{
					title: 'mihomo inside the node',
					text: 'mihomo is embedded as a library, so traffic is counted per user on every connection. Plans by time and gigabytes, device limits and device binding.',
					href: docs('en', 'start/overview/'),
				},
				{
					title: 'Fifteen protocols',
					text: 'VLESS REALITY with Vision, XHTTP and gRPC, Hysteria2, TUIC, AnyTLS, Trojan, TrustTunnel, Mieru and more. Each one comes from a preset with keys generated for you.',
					href: docs('en', 'protocols/presets/'),
				},
				{
					title: 'Survives blocking on its own',
					text: 'When clients from different networks stop reaching a protocol, the panel moves it to a free HTTPS port or picks a new camouflage site next to the server.',
					href: docs('en', 'protocols/autotune/'),
				},
				{
					title: 'Nodes, cascades and WARP',
					text: 'A remote node installs with one command and a key, and updates from the panel. A protocol can leave through another node or through Cloudflare WARP.',
					href: docs('en', 'nodes/nodes/'),
				},
				{
					title: 'A subscription for every app',
					text: 'The subscription recognises Happ, v2RayTun, Koala Clash, Clash Verge, Hiddify and others and gives each only what it runs. Clash apps get routing by service.',
					href: docs('en', 'subscriptions/page/'),
				},
				{
					title: 'Telegram bot and Mini App',
					text: 'Subscribers check their term and devices, buy and renew. Payments through Telegram Stars and marketplace adapters: YooKassa, CryptoBot.',
					href: docs('en', 'billing/telegram/'),
				},
				{
					title: 'Promo codes, pools and packages',
					text: 'Discounts, bonus days and traffic. A separate limit for some protocols, extra gigabytes sold in the bot, a free trial.',
					href: docs('en', 'billing/promo-codes/'),
				},
				{
					title: 'Moving from other panels',
					text: 'Import users from Marzban, PasarGuard and Remnawave. Old subscription links keep working once the old domain points to mikan.',
					href: docs('en', 'operations/import/'),
				},
				{
					title: 'Calm updates',
					text: 'PostgreSQL, signed releases, a backup before every update and a rollback if the new version does not come up. A database backup can arrive in Telegram daily.',
					href: docs('en', 'operations/updates/'),
				},
			],
		},
		protocols: {
			title: 'Every app gets what it can run',
			lead: 'The subscription recognises the app and never hands it a protocol it cannot start, so nobody is left with a key that does not work.',
			headers: ['Protocol', 'Clash apps', 'Happ, v2RayTun', 'Hiddify, Karing'],
			yes: 'yes',
			no: 'no',
			note: '¹ One key for all users: no per-user accounting or limits on these.',
			more: 'All presets and their settings',
			moreHref: docs('en', 'protocols/presets/'),
		},
		phones: {
			title: 'Subscription page and Mini App',
			text: 'A subscriber opens the link in a browser or in Telegram and sees the term, traffic and devices. The platform is detected, and the Add button opens the subscription right in the app: Happ, v2RayTun, Koala Clash, Hiddify and more.',
			link: 'About subscriptions and apps',
			href: docs('en', 'subscriptions/page/'),
			alts: ['The subscription page on a phone', 'The panel overview on a phone'],
		},
		roadmap: {
			badge: 'Free core roadmap',
			title: 'What we are building in Mikan 2.0',
			lead: 'These features are planned for the free core. Development is not gated by fundraising targets.',
			all: 'Full roadmap',
		},
		cta: {
			title: 'Mikan 2.0 — built in the open',
			text: 'The one-command supported installer will be available after the first signed release. Until then, explore the source and roadmap.',
			docs: 'Installation guide',
			sponsorTitle: 'No features hidden behind fundraising',
			sponsorText: 'Core roadmap features are planned to be free. Status and technical limits are published openly.',
			sponsor: 'Open the roadmap',
			goals: 'Development plan',
			github: 'Source code on GitHub',
		},
		header: {
			label: 'Project',
			githubLabel: 'GitHub: the mikan source code',
			support: 'Roadmap',
			supportLabel: 'Mikan 2.0 roadmap',
			otherLanguage: 'Русская версия',
		},
		footer: {
			label: 'Project links',
			license: 'mikan is free software under the GNU GPL v3. It embeds mihomo (GPL-3.0).',
			releases: 'Releases',
			issues: 'Report a problem',
			sponsor: 'Sponsor',
		},
	},
} as const;

// From README.md: Clash apps (mihomo core), Xray apps, sing-box apps.
export const protocolMatrix: [string, boolean, boolean, boolean, boolean?][] = [
	['VLESS · REALITY · Vision', true, true, true],
	['VLESS · REALITY · XHTTP', true, true, false],
	['VLESS · REALITY · gRPC', true, true, true],
	['VLESS PQ', true, true, false],
	['Trojan · REALITY', true, true, true],
	['Hysteria2', true, true, true],
	['TUIC v5', true, false, true],
	['AnyTLS', true, false, true],
	['TrustTunnel', true, false, false],
	['ShadowQUIC', true, false, false],
	['Mieru', true, false, false],
	['Shadowsocks-2022', true, true, true, true],
	['Sudoku', true, false, false, true],
	['Snell', true, false, false, true],
	['VMess', true, true, true],
];
