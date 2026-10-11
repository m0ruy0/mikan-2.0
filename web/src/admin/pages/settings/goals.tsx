import { t } from "../../../i18n";

/**
 * Mikan 2.0's public roadmap. These are planned features of the free core,
 * not crowdfunding goals; this card intentionally does not read release goals
 * or expose donation links.
 */
export function GoalsCard() {
  const features = [
    {
      id: "pterodactyl",
      title: t("goals.features.pterodactyl.title"),
      description: t("goals.features.pterodactyl.description"),
    },
    {
      id: "migration",
      title: t("goals.features.migration.title"),
      description: t("goals.features.migration.description"),
    },
    {
      id: "external",
      title: t("goals.features.external.title"),
      description: t("goals.features.external.description"),
    },
    {
      id: "adminBot",
      title: t("goals.features.adminBot.title"),
      description: t("goals.features.adminBot.description"),
    },
    {
      id: "addons",
      title: t("goals.features.addons.title"),
      description: t("goals.features.addons.description"),
    },
    {
      id: "api",
      title: t("goals.features.api.title"),
      description: t("goals.features.api.description"),
    },
  ];

  return (
    <section className="card glass reveal" aria-labelledby="goals-title">
      <div className="card-head">
        <div>
          <h2 className="card-title" id="goals-title">
            {t("goals.title")}
          </h2>
          <div className="card-sub">{t("goals.sub")}</div>
        </div>
      </div>
      <ul className="row-list">
        {features.map((feature) => (
          <li key={feature.id} className="py-3">
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0">
                <div className="text-[13px] font-semibold text-[var(--ink-900)]">
                  {feature.title}
                </div>
                <p className="mt-1 text-xs text-[var(--ink-500)]">
                  {feature.description}
                </p>
              </div>
              <span className="shrink-0 rounded-full border border-[var(--hairline)] px-2 py-1 text-[11px] text-[var(--ink-500)]">
                {t("goals.planned")}
              </span>
            </div>
            <div className="mt-2 text-xs font-medium text-[var(--ink-600)]">
              {t("goals.freeCore")}
            </div>
          </li>
        ))}
      </ul>
    </section>
  );
}
