import { t } from "../../../i18n";

/**
 * Mikan 2.0's public roadmap. These are planned features of the free core,
 * not crowdfunding goals; this card intentionally does not read release goals
 * or expose donation links.
 */
const FEATURES = [
  ["pterodactyl", "goals.features.pterodactyl"],
  ["migration", "goals.features.migration"],
  ["external", "goals.features.external"],
  ["adminBot", "goals.features.adminBot"],
  ["addons", "goals.features.addons"],
  ["api", "goals.features.api"],
] as const;

export function GoalsCard() {
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
        {FEATURES.map(([id, key]) => (
          <li key={id} className="py-3">
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0">
                <div className="text-[13px] font-semibold text-[var(--ink-900)]">
                  {t(`${key}.title`)}
                </div>
                <p className="mt-1 text-xs text-[var(--ink-500)]">
                  {t(`${key}.description`)}
                </p>
              </div>
              <span className="shrink-0 rounded-full border border-[var(--line)] px-2 py-1 text-[11px] text-[var(--ink-500)]">
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
