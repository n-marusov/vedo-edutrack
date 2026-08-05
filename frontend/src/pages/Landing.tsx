import {
  AlertCircle,
  ArrowRight,
  BadgeCheck,
  Briefcase,
  Check,
  Clock,
  Compass,
  Eye,
  FlaskConical,
  FolderKanban,
  HelpCircle,
  Info,
  Layers,
  Link2,
  Map as MapIcon,
  Minus,
  Navigation,
  Play,
  Plus,
  Quote,
  RefreshCw,
  Search,
  Send,
  Sparkles,
  Target,
  Unlink,
  Users,
} from 'lucide-react';
import { type ReactNode, useEffect, useState } from 'react';
import '../styles/pixso-variables.css';

const THEME_KEY = 'pixso-theme';

function getInitialTheme(): string {
  const stored = localStorage.getItem(THEME_KEY);
  if (stored === 'dark' || stored === 'light') return stored;
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function useTheme() {
  const [theme, setTheme] = useState<string>(() => getInitialTheme());

  useEffect(() => {
    document.documentElement.setAttribute('data-collection-3-4-mode', theme);
    localStorage.setItem(THEME_KEY, theme);
  }, [theme]);

  useEffect(() => {
    const mq = window.matchMedia('(prefers-color-scheme: dark)');
    const handler = (e: MediaQueryListEvent) => {
      const next = e.matches ? 'dark' : 'light';
      setTheme(next);
      document.documentElement.setAttribute('data-collection-3-4-mode', next);
      localStorage.setItem(THEME_KEY, next);
    };
    mq.addEventListener('change', handler);
    return () => mq.removeEventListener('change', handler);
  }, []);

  return { theme, setTheme };
}

// ---------- REUSABLE UI ----------

function SectionBadge({ children }: { children: ReactNode }) {
  return (
    <div className="inline-flex items-center gap-1.5 rounded-full bg-[var(--primary-warm)] px-3.5 py-1.5">
      <div className="size-1.5 rounded-full bg-[var(--accent-foreground)]" />
      <span className="text-[11px] font-bold tracking-[2px] text-[var(--accent-foreground)]">
        {children}
      </span>
    </div>
  );
}

function SectionHeader({
  badge,
  title,
  subtitle,
}: {
  badge: string;
  title: string;
  subtitle?: string;
}) {
  return (
    <div className="flex w-full flex-col items-center gap-4">
      <SectionBadge>{badge}</SectionBadge>
      <h2 className="text-center text-[40px] font-extrabold leading-[110%] tracking-0 text-[var(--foreground)] lg:text-[48px]">
        {title}
      </h2>
      {subtitle && (
        <p className="max-w-[720px] text-center text-[18px] leading-[160%] text-[var(--muted-foreground)] lg:text-[20px]">
          {subtitle}
        </p>
      )}
    </div>
  );
}

function OrangeButton({
  children,
  className = '',
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <button
      type="button"
      className={`inline-flex items-center gap-3 rounded-full bg-[var(--primary-warm)] px-6 py-3.5 text-[16px] font-semibold text-[var(--accent-foreground)] transition-opacity hover:opacity-90 ${className}`}
    >
      {children}
    </button>
  );
}

function Card({
  children,
  className = '',
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <div
      className={`flex w-full flex-1 flex-col rounded-[28px] bg-[var(--card)] shadow-[0_12px_48px_-16px_rgba(124,45,18,0.149)] ${className}`}
    >
      {children}
    </div>
  );
}

function ArrowIconCircle() {
  return (
    <div className="flex size-8 items-center justify-center rounded-full bg-[var(--accent-foreground)]">
      <ArrowRight className="size-4 text-[var(--primary-foreground)]" strokeWidth={2.5} />
    </div>
  );
}

// ---------- ASSET ICONS (custom inline SVGs) ----------

function IconHand() {
  return (
    <svg width="28" height="28" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
      <title>Hand icon</title>
      <path
        d="M13.9999 2.33337C16.5772 2.33337 18.6666 4.42271 18.6666 7.00004V11.6667C18.6666 14.244 16.5772 16.3334 13.9999 16.3334C11.4226 16.3334 9.33325 14.244 9.33325 11.6667V7.00004C9.33325 4.42271 11.4226 2.33337 13.9999 2.33337Z"
        fill="var(--primary-warm)"
      />
      <path
        d="M7.00008 12.8334C7.00008 12.8334 5.83341 14.4667 5.83341 16.3334C5.83341 18.8427 7.88275 20.8334 10.5001 20.8334H17.5001C20.1174 20.8334 22.1667 18.8427 22.1667 16.3334C22.1667 14.4667 21.0001 12.8334 21.0001 12.8334"
        stroke="var(--primary-warm)"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      <path
        d="M10.5 20.8334V23.3334M14 20.8334V25.6667M17.5 20.8334V23.3334"
        stroke="var(--primary-warm)"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
    </svg>
  );
}

// ---------- SECTION: HEADER ----------

const NAV_LINKS = ['Как это работает', 'Возможности', 'Тарифы', 'Отзывы', 'FAQ'];

function Header() {
  return (
    <header className="sticky top-0 z-50 w-full px-5 py-4 lg:px-20">
      <div className="flex w-full items-center justify-between rounded-full border border-[var(--border)] bg-[var(--card)] px-4 py-2.5 shadow-[0_12px_48px_-16px_rgba(124,45,18,0.149)] lg:px-8">
        <a href="/" className="flex items-center gap-3">
          <IconHand />
          <span className="text-[22px] font-bold tracking-[-0.5px] text-[var(--foreground)]">
            Дай пять
          </span>
        </a>
        <nav className="hidden items-center gap-8 lg:flex">
          {NAV_LINKS.map((link) => (
            <a
              key={link}
              href="/"
              className="text-[16px] font-medium text-[var(--muted-foreground)] transition-colors hover:text-[var(--foreground)]"
            >
              {link}
            </a>
          ))}
        </nav>
        <div className="flex items-center gap-5">
          <div className="hidden items-center gap-2 md:flex">
            <Briefcase className="size-4 text-[var(--muted-foreground)]" strokeWidth={1.5} />
            <span className="text-[16px] text-[var(--muted-foreground)]">Для бизнеса</span>
          </div>
          <div className="hidden items-center gap-2 md:flex">
            <Layers className="size-4 text-[var(--muted-foreground)]" strokeWidth={1.5} />
            <span className="text-[16px] text-[var(--muted-foreground)]">VEDO Hub</span>
          </div>
          <OrangeButton className="hidden px-5 py-2.5 sm:inline-flex">
            Попробовать бесплатно
          </OrangeButton>
        </div>
      </div>
    </header>
  );
}

// ---------- SECTION: HERO ----------

function Hero() {
  return (
    <section className="flex flex-col items-center justify-between gap-12 px-5 pb-20 pt-16 lg:flex-row lg:px-20 lg:pt-[100px]">
      <div className="flex max-w-[720px] flex-col items-start gap-8">
        <div className="inline-flex items-center gap-2 rounded-full bg-[var(--muted)] px-4 py-2">
          <Sparkles className="size-4 text-[var(--accent-warm)]" strokeWidth={2} />
          <span className="text-[14px] font-semibold text-[var(--muted-foreground)]">
            Не зубрёжка, а понимание
          </span>
        </div>
        <h1 className="text-[52px] font-extrabold leading-[110%] tracking-[-2px] text-[var(--foreground)] lg:text-[72px]">
          Ребёнок учится думать.
          <br />
          Вы — направлять.
        </h1>
        <p className="max-w-[640px] text-[20px] leading-[160%] text-[var(--muted-foreground)]">
          Платформа для построения индивидуальных маршрутов обучения на основе графа знаний,
          диагностики пробелов, отслеживания прогресса и ИИ-ассистента
        </p>
        <div className="flex flex-wrap items-center gap-4">
          <OrangeButton className="gap-4 px-6 py-4 text-[17px]">
            Построить маршрут за 5 минут
            <ArrowIconCircle />
          </OrangeButton>
          <button
            type="button"
            className="rounded-full border border-[var(--border)] px-7 py-4 text-[17px] font-semibold text-[var(--foreground)] transition-colors hover:bg-[var(--muted)]"
          >
            Узнать больше
          </button>
        </div>
      </div>

      {/* Hero illustration card */}
      <div className="relative w-[640px] h-[580px] shrink-0 rounded-[40px] bg-gradient-to-br from-[#FFF7ED] via-[#FFEDD5] to-[#ECFEFF] max-lg:hidden">
        <div className="absolute left-[60px] top-[50px] flex w-[520px] h-[480px] flex-col gap-5 rounded-[32px] bg-[var(--card)] p-8 shadow-[0_16px_64px_-20px_rgba(124,45,18,0.149)]">
          {/* Card header */}
          <div className="flex w-full items-center justify-between">
            <span className="text-[20px] font-bold text-[var(--foreground)]">Карта знаний</span>
            <div className="inline-flex items-center gap-1.5 rounded-full bg-[var(--accent-warm)] px-3 py-1.5">
              <Check className="size-3 text-[var(--primary-foreground)]" strokeWidth={3} />
              <span className="text-[12px] font-bold text-[var(--primary-foreground)]">ФГОС</span>
            </div>
          </div>
          {/* Progress bar */}
          <div className="flex w-full flex-col gap-3 rounded-[20px] bg-[var(--muted)] p-5">
            <span className="text-[16px] font-semibold text-[var(--foreground)]">
              Алгебра: деление → дроби
            </span>
            <div className="h-2.5 w-full rounded-full bg-[var(--border)]">
              <div className="h-2.5 w-[300px] rounded-full bg-[var(--primary-warm)]" />
            </div>
            <div className="flex w-full justify-between">
              <span className="text-[13px] text-[var(--muted-foreground)]">Пройдено 75%</span>
              <span className="text-[13px] font-semibold text-[var(--primary-warm)]">
                12/16 тем
              </span>
            </div>
          </div>
          {/* Gap alert */}
          <div className="flex w-full items-center gap-3 rounded-[20px] border border-[var(--border)] bg-[#FFF7ED] p-4">
            <AlertCircle className="size-5 shrink-0 text-[var(--primary-warm)]" strokeWidth={1.5} />
            <div className="flex flex-col">
              <span className="text-[15px] font-bold text-[var(--foreground)]">
                Пробел: Проценты
              </span>
              <span className="text-[13px] text-[var(--muted-foreground)]">
                Нужны деление → дроби
              </span>
            </div>
          </div>
          {/* Route steps */}
          <div className="flex w-full flex-col gap-2.5 rounded-[20px] bg-[var(--muted)] p-4">
            <span className="text-[15px] font-semibold text-[var(--foreground)]">
              Маршрут до аттестации
            </span>
            <div className="flex w-full items-center gap-2">
              <div className="flex h-9 items-center justify-center rounded-xl bg-[var(--primary-warm)] px-4">
                <span className="text-[13px] font-semibold text-[var(--primary-foreground)]">
                  Деление
                </span>
              </div>
              <ArrowRight
                className="size-4 shrink-0 text-[var(--muted-foreground)]"
                strokeWidth={1.5}
              />
              <div className="flex h-9 items-center justify-center rounded-xl bg-[var(--secondary-calm)] px-4">
                <span className="text-[13px] font-semibold text-[var(--primary-foreground)]">
                  Дроби
                </span>
              </div>
              <ArrowRight
                className="size-4 shrink-0 text-[var(--muted-foreground)]"
                strokeWidth={1.5}
              />
              <div className="flex h-9 items-center justify-center rounded-xl bg-[var(--accent-warm)] px-4">
                <span className="text-[13px] font-semibold text-[var(--primary-foreground)]">
                  Проценты
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

// ---------- SECTION: PROBLEM ----------

const PROBLEMS = [
  {
    icon: AlertCircle,
    title: 'Программа перегружена',
    desc: 'Школьную программу невозможно осмысленно пройти в отведённое время: либо поверхностно, либо не успеваем.',
  },
  {
    icon: Search,
    title: 'Пробелы в знаниях',
    desc: 'В декабре выяснилось, что ребёнок не знает тему, которой не было в купленных курсах. Где ещё пробелы?',
  },
  {
    icon: Unlink,
    title: 'Предметы изолированы',
    desc: 'Физика отдельно от математики, биология — от химии. Ребёнок не видит связей между явлениями.',
  },
  {
    icon: Clock,
    title: 'Единый темп не работает',
    desc: 'По математике — 7 класс, по биологии — 8. Это нормально. Но как это учесть в плане?',
  },
  {
    icon: HelpCircle,
    title: 'Непонятно зачем',
    desc: 'Ребёнок спрашивает: «Зачем мне это?» — а ответить нечего. Каждое знание должно иметь смысл.',
  },
];

function ProblemSection() {
  return (
    <section className="flex w-full flex-col bg-[var(--background)] px-5 py-[100px] lg:px-20 lg:py-[140px]">
      <div className="mx-auto flex w-full max-w-[1760px] flex-col items-center gap-[72px]">
        <SectionHeader
          badge="ПРОБЛЕМА"
          title="Это про вас?"
          subtitle="Знакомые ситуации для родителей на семейном образовании"
        />
        <div className="grid w-full grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5">
          {PROBLEMS.map((item) => (
            <Card key={item.title}>
              <div className="flex flex-col gap-4 p-8">
                <div className="flex size-8 items-center justify-center rounded-lg bg-[var(--muted)]">
                  <item.icon className="size-4 text-[var(--primary-warm)]" strokeWidth={2} />
                </div>
                <h3 className="text-[20px] font-bold leading-[130%] text-[var(--foreground)]">
                  {item.title}
                </h3>
                <p className="w-full text-[16px] leading-[160%] text-[var(--muted-foreground)]">
                  {item.desc}
                </p>
              </div>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}

// ---------- SECTION: BRIDGE ----------

function Bridge() {
  return (
    <section className="flex w-full flex-col items-center bg-[var(--muted)] px-5 py-20 lg:px-20">
      <p className="text-center text-[28px] font-semibold text-[var(--muted-foreground)] lg:text-[32px]">
        Так было.
      </p>
      <p className="text-center text-[28px] font-bold text-[var(--foreground)] lg:text-[32px]">
        А теперь — так
      </p>
    </section>
  );
}

// ---------- SECTION: SOLUTION ----------

const STEPS = [
  {
    num: '1',
    title: 'Задаём цели и ограничения',
    desc: 'Цели к аттестации, время, нагрузка, темп. ИИ-ассистент поможет сформулировать и скорректировать.',
    example: 'Пример: «К маю — аттестация, по алгебре подтянуть до «4»»',
  },
  {
    num: '2',
    title: 'Подключаем карту знаний',
    desc: 'Подключаем готовую карту от методистов и сообщества: темы, связи, актуальный ФГОС.',
    example: 'Пример: «Проценты → биология → химия → география»',
  },
  {
    num: '3',
    title: 'Находим пробелы',
    desc: 'Система показывает, какие темы пропущены и почему это важно.',
    example: 'Пример: «Деление → Дроби → Проценты. Без деления — пробел»',
  },
  {
    num: '4',
    title: 'Видим прогноз до аттестации',
    desc: 'Видите, что успеете до экзамена. Планируете с уверенностью.',
    example: 'Пример: «До ОГЭ 8 месяцев. Успеете 95% программы»',
  },
];

const CONCEPTS = [
  {
    icon: RefreshCw,
    title: 'Спиральное обучение',
    desc: 'Тема «Функции» возвращается в 7, 8, 9 классах — с возрастающей глубиной. Маршрут строит витки, а не линейную цепочку.',
    selected: true,
  },
  {
    icon: FolderKanban,
    title: 'Проектное погружение',
    desc: 'Модули из разных предметов группируются вокруг проекта «Экосистема»: биология + математика + химия + география.',
    selected: false,
  },
  {
    icon: FlaskConical,
    title: 'Свободное обучение',
    desc: 'Ребёнок идёт от интереса: динозавры → биология → география → химия. Маршрут строится вокруг увлечений, а не наоборот.',
    selected: false,
  },
];

function SolutionSection() {
  return (
    <section className="flex w-full flex-col bg-[var(--background)] px-5 py-[100px] lg:px-20 lg:py-[140px]">
      <div className="mx-auto flex w-full max-w-[1760px] flex-col gap-[72px]">
        <SectionHeader
          badge="РЕШЕНИЕ"
          title="Как работает «Дай пять»"
          subtitle="Представьте, что все школьные темы — это города на карте, а связи между ними — дороги. Мы ведём вас по маршруту."
        />
        <div className="flex flex-wrap items-stretch gap-8">
          {STEPS.map((step) => (
            <Card key={step.num} className="min-w-[280px] flex-1">
              <div className="flex flex-col gap-6 p-8 lg:p-10">
                <div className="flex size-16 items-center justify-center rounded-full bg-[var(--primary-warm)]">
                  <span className="text-[32px] font-extrabold text-[var(--accent-foreground)]">
                    {step.num}
                  </span>
                </div>
                <h3 className="text-[22px] font-bold leading-[130%] text-[var(--foreground)] lg:text-[24px]">
                  {step.title}
                </h3>
                <p className="w-full text-[16px] leading-[160%] text-[var(--muted-foreground)]">
                  {step.desc}
                </p>
                <p className="text-[14px] font-medium leading-[150%] text-[var(--primary-warm)]">
                  {step.example}
                </p>
              </div>
            </Card>
          ))}
        </div>

        {/* Concept picker */}
        <div className="flex w-full flex-col items-center gap-8">
          <div className="flex flex-col items-center gap-3">
            <h3 className="text-center text-[28px] font-bold text-[var(--foreground)]">
              Не знаете, с чего начать?
            </h3>
            <p className="text-center text-[17px] font-medium leading-[150%] text-[var(--muted-foreground)]">
              Выберите педагогическую концепцию
            </p>
          </div>
          <div className="flex w-full flex-wrap items-stretch gap-6">
            {CONCEPTS.map((item) => (
              <div
                key={item.title}
                className="flex flex-1 flex-col gap-3.5 rounded-[24px] border-2 border-[var(--border)] bg-[var(--card)] p-6 shadow-[0_12px_48px_-16px_rgba(124,45,18,0.149)] lg:p-7"
              >
                <div className="flex size-10 items-center justify-center rounded-xl bg-[var(--muted)]">
                  <item.icon className="size-5 text-[var(--primary-warm)]" strokeWidth={2} />
                </div>
                <h4 className="text-[20px] font-bold text-[var(--foreground)]">{item.title}</h4>
                <p className="text-[15px] leading-[160%] text-[var(--muted-foreground)]">
                  {item.desc}
                </p>
                <div
                  className={`inline-flex w-fit items-center gap-2 rounded-full px-3.5 py-2 ${
                    item.selected ? 'bg-[var(--primary-warm)]' : 'bg-[var(--muted)]'
                  }`}
                >
                  {item.selected && (
                    <Check className="size-3.5 text-[var(--primary-foreground)]" strokeWidth={3} />
                  )}
                  <span
                    className={`text-[13px] font-semibold ${
                      item.selected
                        ? 'text-[var(--primary-foreground)]'
                        : 'text-[var(--muted-foreground)]'
                    }`}
                  >
                    {item.selected ? 'Выбрано' : 'Выбрать'}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}

// ---------- SECTION: BENEFITS ----------

const BENEFITS_ROW_1 = [
  {
    icon: Eye,
    title: 'Вижу пробелы до аттестации',
    desc: 'Система заранее покажет, какие темы не пройдены. Не будет сюрпризов на экзамене.',
    instead: 'Вместо: «Надеюсь, всё успеем»',
  },
  {
    icon: Clock,
    title: 'Учу в своём темпе',
    desc: 'По математике — 7 класс, по биологии — 8. Это нормально. Система учитывает разные уровни.',
    instead: 'Вместо: «Все должны быть на одном уровне»',
  },
  {
    icon: Users,
    title: 'Одна панель для всех детей',
    desc: 'Трое детей? Не проблема. Видите прогресс каждого в одном месте.',
    instead: 'Вместо: «Три таблицы Excel и два чата»',
  },
];

const BENEFITS_ROW_2 = [
  {
    icon: Link2,
    title: 'Знаю, зачем это учить',
    desc: 'Каждая тема связана с реальной жизнью и другими предметами. Ребёнок видит смысл.',
    instead: 'Вместо: «Учи, потому что так надо»',
  },
  {
    icon: MapIcon,
    title: 'Карта знаний показывает связи',
    desc: 'Визуальная карта с темами и связями между предметами. Физика не отдельно от математики.',
    instead: 'Вместо: «Предметы как отдельные острова»',
  },
  {
    icon: BadgeCheck,
    title: 'Успеваем подготовиться к аттестации',
    desc: 'Система показывает покрытие стандарта: какие темы закрыты, какие остались до аттестации. Без ручной сверки в Excel.',
    instead: 'Вместо: «Ручная сверка с ФГОС в Excel»',
  },
];

function BenefitCardList({ items }: { items: typeof BENEFITS_ROW_1 }) {
  return (
    <div className="flex w-full flex-wrap items-stretch gap-6">
      {items.map((item) => (
        <Card key={item.title} className="min-w-[280px] flex-1">
          <div className="flex flex-col gap-5 p-8 lg:p-10">
            <div className="flex size-10 items-center justify-center rounded-xl bg-[var(--muted)]">
              <item.icon className="size-5 text-[var(--secondary-calm)]" strokeWidth={2} />
            </div>
            <h3 className="text-[20px] font-bold leading-[130%] text-[var(--foreground)] lg:text-[22px]">
              {item.title}
            </h3>
            <p className="w-full text-[16px] leading-[160%] text-[var(--muted-foreground)]">
              {item.desc}
            </p>
            <p className="text-[14px] font-medium leading-[150%] text-[var(--secondary-calm)]">
              {item.instead}
            </p>
          </div>
        </Card>
      ))}
    </div>
  );
}

function BenefitsSection() {
  return (
    <section className="flex w-full flex-col bg-[var(--muted)] px-5 py-[100px] lg:px-20 lg:py-[140px]">
      <div className="mx-auto flex w-full max-w-[1760px] flex-col gap-[72px]">
        <SectionHeader
          badge="ПРЕИМУЩЕСТВА"
          title="Почему родители выбирают «Дай пять»"
          subtitle="Вместо метода проб и ошибок — система, которая работает"
        />
        <div className="flex w-full flex-col gap-8">
          <BenefitCardList items={BENEFITS_ROW_1} />
          <BenefitCardList items={BENEFITS_ROW_2} />
        </div>
        <div className="flex w-full items-center justify-center gap-4 rounded-xl bg-[var(--card)] px-8 py-8 lg:px-12">
          <Target className="size-6 shrink-0 text-[var(--primary-warm)]" strokeWidth={2} />
          <p className="text-center text-[18px] font-medium leading-[150%] text-[var(--foreground)]">
            Почему это работает: карта знаний основана на принципе межпредметных связей — каждая
            тема имеет понятные «зачем» и «где пригодится». А ИИ-ассистент ответит на вопросы и
            подскажет, как лучше выстроить маршрут.
          </p>
        </div>
      </div>
    </section>
  );
}

// ---------- SECTION: PRICING ----------

const PLANS = [
  {
    icon: MapIcon,
    name: 'Карта',
    price: '0 ₽',
    period: '/ 7 дней',
    description: 'Диагностика пробелов и карта тем на 7 дней.',
    features: [
      { text: 'Карта всех тем по предметам', highlight: false },
      { text: 'Диагностика пробелов до аттестации', highlight: false },
      { text: '1 ребёнок в профиле', highlight: false },
      { text: 'Без привязки банковской карты', highlight: false },
      { text: 'Истории «зачем это знать»', highlight: false },
      { text: 'Маршрут и расписание', highlight: false },
      { text: 'Планирование бюджета и ресурсов', highlight: false },
    ],
    cta: 'Попробовать 7 дней',
  },
  {
    icon: Compass,
    name: 'Компас',
    price: '490 ₽',
    period: '/ месяц',
    annual: { price: '4 410 ₽', oldPrice: '5 880 ₽', savings: 'Экономия 1 470 ₽' },
    description: 'Направление для всей семьи. Истории, маршрут и расписание, 30 запросов к ИИ.',
    features: [
      { text: 'Всё из Карты', highlight: false },
      { text: 'Истории и контекст «зачем это знать»', highlight: true },
      { text: 'Маршрут и расписание: что доступно сейчас, что дальше, где цель', highlight: true },
      { text: 'Полный учебный маршрут', highlight: false },
      { text: 'Отслеживание прогресса', highlight: false },
      { text: 'До 2 детей в одном профиле', highlight: false },
      { text: 'ИИ-ассистент (30 запросов/мес)', highlight: false },
    ],
    cta: 'Выбрать маршрут',
  },
  {
    icon: Navigation,
    name: 'Навигатор',
    badge: 'Самый популярный',
    price: '990 ₽',
    period: '/ месяц',
    annual: { price: '8 910 ₽', oldPrice: '11 880 ₽', savings: 'Экономия 2 970 ₽' },
    description:
      'Полная система для осмысленной подготовки. 300 запросов к ИИ, ресурсы, проекты и отчёты.',
    features: [
      { text: 'Всё из Компаса', highlight: false },
      { text: 'Без ограничений по количеству детей', highlight: false },
      { text: 'ИИ-ассистент по планированию (300 запросов/мес)', highlight: true },
      {
        text: 'Планирование бюджета и ресурсов — сколько денег и времени до аттестации',
        highlight: true,
      },
      { text: 'Проектные идеи на стыках предметов', highlight: false },
      { text: 'Экспорт отчётов и сертификатов', highlight: false },
      { text: 'Приоритетная поддержка', highlight: false },
    ],
    cta: 'Получить максимум',
  },
];

function PricingCard({ plan }: { plan: (typeof PLANS)[number] }) {
  return (
    <div className="flex w-full flex-1 flex-col rounded-[24px] bg-[var(--card)] p-10 shadow-[0_12px_48px_-16px_rgba(124,45,18,0.149)]">
      {/* Title row */}
      <div className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-2.5">
          <plan.icon className="size-7 text-[var(--primary-warm)]" strokeWidth={2} />
          <h3 className="text-[28px] font-bold text-[var(--foreground)]">{plan.name}</h3>
        </div>
        {plan.badge && (
          <span className="inline-flex items-center rounded-full bg-[var(--primary-warm)] px-3.5 py-1.5 text-[13px] font-semibold text-[var(--accent-foreground)]">
            {plan.badge}
          </span>
        )}
      </div>

      {/* Price row */}
      <div className="mt-4 flex items-end gap-1">
        <span className="text-[48px] font-extrabold leading-none text-[var(--foreground)]">
          {plan.price}
        </span>
        <span className="pb-1 text-[18px] text-[var(--muted-foreground)]">{plan.period}</span>
      </div>

      {/* Annual row */}
      {plan.annual && (
        <div className="mt-3 flex items-center gap-2.5">
          <span className="text-[16px] font-bold text-[var(--foreground)]">
            {plan.annual.price}
          </span>
          <span className="text-[14px] text-[var(--muted-foreground)] line-through">
            {plan.annual.oldPrice}
          </span>
          <span className="inline-flex items-center rounded-full bg-[rgba(13,148,136,0.12)] px-2.5 py-1 text-[13px] font-semibold text-[rgba(15,118,110,1)]">
            {plan.annual.savings}
          </span>
        </div>
      )}

      {/* Description */}
      <p className="mt-4 text-[16px] leading-[150%] text-[var(--muted-foreground)]">
        {plan.description}
      </p>

      {/* Divider */}
      <div className="mt-5 h-px w-full bg-[var(--border)]" />

      {/* Features */}
      <ul className="mt-5 flex flex-col gap-3">
        {plan.features.map((feature) => (
          <li key={feature.text}>
            {feature.highlight ? (
              <div className="flex items-start gap-3 rounded-xl bg-[var(--muted)] px-3 py-2.5">
                <Check
                  className="mt-0.5 size-[18px] shrink-0 text-[var(--secondary-calm)]"
                  strokeWidth={2.5}
                />
                <span className="text-[16px] font-semibold leading-[150%] text-[var(--foreground)]">
                  {feature.text}
                </span>
              </div>
            ) : (
              <div className="flex items-start gap-3">
                <Check
                  className="mt-0.5 size-[18px] shrink-0 text-[var(--secondary-calm)]"
                  strokeWidth={2.5}
                />
                <span className="text-[16px] leading-[150%] text-[var(--muted-foreground)]">
                  {feature.text}
                </span>
              </div>
            )}
          </li>
        ))}
      </ul>

      {/* CTA */}
      <button
        type="button"
        className="mt-8 flex w-full items-center justify-center gap-3 rounded-full bg-[var(--primary-warm)] py-4 text-[16px] font-semibold text-[var(--accent-foreground)] transition-opacity hover:opacity-90"
      >
        {plan.cta}
        <span className="flex size-8 items-center justify-center rounded-full bg-[var(--accent-foreground)]">
          <ArrowRight className="size-4 text-[var(--primary-foreground)]" strokeWidth={2.5} />
        </span>
      </button>
    </div>
  );
}

function PricingSection() {
  return (
    <section className="flex w-full flex-col bg-[var(--background)] px-5 py-[100px] lg:px-20 lg:py-[140px]">
      <div className="mx-auto flex w-full max-w-[1760px] flex-col items-center gap-[72px]">
        <SectionHeader
          badge="ТАРИФЫ"
          title="Выберите свой маршрут к аттестации"
          subtitle="Диагностика бесплатно. Мотивация — в Компасе. Экономия на репетиторах — в Навигаторе."
        />
        <div className="grid w-full grid-cols-1 gap-8 lg:grid-cols-3">
          {PLANS.map((plan) => (
            <PricingCard key={plan.name} plan={plan} />
          ))}
        </div>
        <div className="flex w-full flex-wrap items-center justify-center gap-3">
          <Info className="size-5 shrink-0 text-[var(--muted-foreground)]" strokeWidth={1.5} />
          <span className="text-center text-[16px] text-[var(--muted-foreground)]">
            Не нашли подходящий вариант? VEDO EduTrack также используют школы, онлайн-платформы и
            корпорации.
          </span>
          <button
            type="button"
            className="inline-flex items-center gap-1 text-[16px] font-semibold text-[var(--primary-warm)] transition-colors hover:underline"
          >
            Перейти в раздел «Для бизнеса»
            <ArrowRight className="size-[18px]" strokeWidth={1.5} />
          </button>
        </div>
      </div>
    </section>
  );
}

// ---------- SECTION: TESTIMONIALS ----------

const TESTIMONIALS = [
  {
    text: '«До «Дай пять» я тратила 3 часа в неделю на составление плана. Сейчас система сама показывает, что учить дальше. Маша стала учиться с удовольствием — видит, зачем ей каждая тема.»',
    author: 'Елена, 38 лет',
    role: 'Мама Маши (6 класс) и Пети (3 класс)',
    result: 'Результат: аттестация сдана на «отлично», пробелов нет',
  },
  {
    text: '«Сын обожает биологию, но алгебра шла туго. «Дай пять» показал: он не понимает дроби → проценты → статистику. Закрыли пробел за 2 месяца. Теперь сам просит задачи по математике!»',
    author: 'Андрей, 42 года',
    role: 'Папа Димы (8 класс)',
    result: 'Результат: оценка по алгебре выросла с 3 до 5',
  },
  {
    text: '«Трое детей на семейном — это хаос. «Дай пять» собрал всё в одну панель: вижу прогресс каждого, дети сами планируют обучение.»',
    author: 'Ольга, 35 лет',
    role: 'Мама троих: Кати (9 класс), Вани (6 класс), Сони (2 класс)',
    result: 'Результат: дети учатся самостоятельно, мама — координатор',
  },
];

const METRICS = [
  { value: '1000+', label: 'тем в базе знаний' },
  { value: '500+', label: 'межпредметных связей' },
  { value: '200+', label: 'семей уже используют' },
];

function TestimonialsSection() {
  return (
    <section className="flex w-full flex-col bg-[var(--muted)] px-5 py-[100px] lg:px-20 lg:py-[140px]">
      <div className="mx-auto flex w-full max-w-[1760px] flex-col items-center gap-[72px]">
        <SectionHeader
          badge="ОТЗЫВЫ"
          title="Реальные истории"
          subtitle="Как «Дай пять» помог семьям на семейном образовании"
        />
        <div className="flex w-full items-stretch gap-10">
          {TESTIMONIALS.map((item) => (
            <div
              key={item.author}
              className="flex w-full flex-1 flex-col gap-6 rounded-[28px] border border-[var(--border)] bg-[var(--card)] p-10 shadow-[0_16px_64px_-20px_rgba(124,45,18,0.149)]"
            >
              <Quote className="size-8 text-[var(--primary-warm)]" strokeWidth={1.5} />
              <p className="text-[18px] leading-[170%] text-[var(--foreground)]">{item.text}</p>
              <div className="flex items-center gap-4 pt-4">
                <div className="flex size-14 shrink-0 items-center justify-center rounded-full bg-[var(--primary-warm)] text-[16px] font-bold text-[var(--accent-foreground)]">
                  {item.author[0]}
                </div>
                <div className="flex flex-col gap-1">
                  <span className="text-[16px] font-bold leading-[130%] text-[var(--foreground)]">
                    {item.author}
                  </span>
                  <span className="text-[14px] leading-[140%] text-[var(--muted-foreground)]">
                    {item.role}
                  </span>
                </div>
              </div>
              <div className="flex w-full items-center gap-3 rounded-[8px] bg-[var(--muted)] px-5 py-4">
                <Check className="size-5 shrink-0 text-[var(--secondary-calm)]" strokeWidth={2.5} />
                <span className="text-[14px] font-semibold leading-[140%] text-[var(--foreground)]">
                  {item.result}
                </span>
              </div>
            </div>
          ))}
        </div>
        <div className="flex w-full items-center justify-center gap-8 rounded-xl bg-[var(--muted)] px-8 py-10">
          {METRICS.map((metric) => (
            <div key={metric.value} className="flex flex-col items-center gap-2">
              <span className="text-[48px] font-extrabold leading-none text-[var(--primary-warm)]">
                {metric.value}
              </span>
              <span className="text-[16px] font-medium leading-[140%] text-[var(--muted-foreground)]">
                {metric.label}
              </span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

// ---------- SECTION: PHILOSOPHY ----------

const PRINCIPLES = [
  {
    icon: Layers,
    title: 'ЦЕЛОСТНОСТЬ',
    subtitle: 'Целостность',
    desc: 'Знания имеют смысл только в системе: одна тема вытекает из другой, предметы переплетаются. Мы строим карту, где эти связи видны, — чтобы ребёнок видел не отдельные куски, а живую картину мира.',
  },
  {
    icon: Compass,
    title: 'СВОБОДА',
    subtitle: 'Свобода без хаоса',
    desc: 'Ребёнок идёт за своим интересом — в рамках маршрута, который ведёт к аттестации. Интерес и ответственность соединяются в едином плане. Мы не загоняем в рамки, но и не отпускаем без карты.',
  },
  {
    icon: Link2,
    title: 'СВЯЗИ',
    subtitle: 'Связь с жизнью',
    desc: 'Каждая тема показывает, где она применяется: в профессии, в быту, в других предметах. Знания не висят в воздухе — они имеют смысл здесь и сейчас.',
  },
  {
    icon: Users,
    title: 'СООБЩЕСТВО',
    subtitle: 'Сила в сообществе',
    desc: 'Карта знаний — живая, её развивают методисты и родители вместе. Это общее дело: вы делитесь опытом, находите единомышленников и влияете на то, как учатся другие.',
    link: 'На платформе VEDO Hub →',
  },
];

function PhilosophySection() {
  return (
    <section className="flex w-full flex-col bg-[var(--background)] px-5 py-[100px] lg:px-20 lg:py-[140px]">
      <div className="mx-auto flex w-full max-w-[1760px] flex-col items-center gap-[72px]">
        <SectionHeader
          badge="ОСНОВА"
          title="Система, которая возвращает логику и смысл"
          subtitle="Мы строим карту знаний, где каждая тема занимает своё место."
        />
        <div className="grid w-full grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-4">
          {PRINCIPLES.map((item) => (
            <div
              key={item.title}
              className="flex w-full flex-col gap-5 rounded-[28px] border border-[var(--border)] bg-[var(--card)] p-10 shadow-[0_16px_64px_-20px_rgba(124,45,18,0.149)]"
            >
              <div className="flex size-8 items-center justify-center">
                <item.icon className="size-8 text-[var(--primary-warm)]" strokeWidth={1.75} />
              </div>
              <span className="text-[13px] font-bold tracking-[1px] text-[var(--muted-foreground)]">
                {item.title}
              </span>
              <h3 className="text-[24px] font-bold leading-[130%] text-[var(--foreground)]">
                {item.subtitle}
              </h3>
              <p className="text-[16px] leading-[160%] text-[var(--muted-foreground)]">
                {item.desc}
              </p>
              {item.link && (
                <span className="mt-auto text-[14px] font-semibold text-[var(--primary-warm)]">
                  {item.link}
                </span>
              )}
            </div>
          ))}
        </div>
        <div className="flex w-full items-center gap-4 rounded-xl bg-[var(--card)] px-12 py-8">
          <Target className="size-6 shrink-0 text-[var(--primary-warm)]" strokeWidth={2} />
          <p className="text-[18px] font-medium leading-[150%] text-[var(--foreground)]">
            «Дай пять» — это как экспедиция. У вас есть карта, где уже проложены маршруты. Но вы
            сами выбираете, по какой тропе идти. Рядом — те, кто прошёл этот путь раньше. И вы
            всегда знаете, где находитесь и что ждёт впереди.
          </p>
        </div>
      </div>
    </section>
  );
}

// ---------- SECTION: FAQ ----------

const FAQS = [
  {
    question: 'Что, если мы не успеем подготовиться к аттестации? Вдруг забудем какую-то тему?',
    answer:
      'Система заранее показывает пробелы: какие темы пройдены, какие остались, и сколько на них нужно времени. Вы видите покрытие ФГОС в реальном времени — ни одна тема не потеряется. Аттестация перестаёт быть стрессом, потому что вы точно знаете: всё под контролем.',
  },
  {
    question: 'Что если я сам не разбираюсь в предмете?',
    answer:
      'Вам не нужно быть методистом. Система показывает, что учить и в каком порядке. Ребёнок учится через пример взрослого — вы организуете, а не преподаёте.',
  },
  {
    question: 'Можно ли вести несколько детей одновременно?',
    answer:
      'Да, одна панель для всех детей. Видите прогресс каждого, управляете планами, получаете уведомления о пробелах.',
  },
  {
    question: 'Это бесплатно?',
    answer:
      'Базовый функционал бесплатен навсегда: карта знаний, диагностика пробелов, план обучения. Продвинутые функции (прогноз аттестации, несколько детей) — по подписке.',
  },
  {
    question: 'Сколько времени займёт настройка?',
    answer:
      '5 минут. Отвечаете на вопросы о ребёнке — система подключает готовую карту знаний. Дальше работаете по плану, который обновляется автоматически.',
  },
  {
    question: 'Мои данные в безопасности?',
    answer:
      'Да, мы соответствуем 152-ФЗ. Данные хранятся в России, шифруются, не передаются третьим лицам. Вы можете удалить аккаунт в любой момент.',
  },
  {
    question: 'Кто создаёт граф знаний?',
    answer:
      'Граф знаний создаётся методистами совместно с родительским сообществом на основе онтологии VEDO Hub. Родители могут предлагать связи между темами, делиться опытом и учебными находками. Это живая система: граф постоянно пополняется и уточняется — академическая точность методистов встречается с реальным родительским опытом.',
  },
];

function FAQItem({ item, index }: { item: (typeof FAQS)[number]; index: number }) {
  const [open, setOpen] = useState(index === 0);
  return (
    <div className="w-full rounded-[28px] bg-[var(--card)] p-8 shadow-[0_12px_48px_-16px_rgba(124,45,18,0.149)]">
      <button
        type="button"
        onClick={() => setOpen((prev) => !prev)}
        className="flex w-full items-center justify-between gap-4 text-left"
        aria-expanded={open}
      >
        <span className="text-[18px] font-bold leading-[140%] text-[var(--foreground)] lg:text-[20px]">
          {item.question}
        </span>
        <div className="flex size-8 shrink-0 items-center justify-center rounded-full bg-[var(--muted)]">
          {open ? (
            <Minus className="size-4 text-[var(--primary-warm)]" strokeWidth={2.5} />
          ) : (
            <Plus className="size-4 text-[var(--primary-warm)]" strokeWidth={2.5} />
          )}
        </div>
      </button>
      {open && (
        <p className="mt-4 text-[18px] leading-[160%] text-[var(--muted-foreground)]">
          {item.answer}
        </p>
      )}
    </div>
  );
}

function FaqSection() {
  return (
    <section className="flex w-full flex-col bg-[var(--muted)] px-5 py-[100px] lg:px-20 lg:py-[140px]">
      <div className="mx-auto flex w-full max-w-[1760px] flex-col items-center gap-[72px]">
        <SectionHeader
          badge="FAQ"
          title="Частые вопросы"
          subtitle="Отвечаем на сомнения до того, как они возникнут"
        />
        <div className="flex w-full flex-col items-center gap-4">
          {FAQS.map((item, index) => (
            <FAQItem key={item.question} item={item} index={index} />
          ))}
        </div>
      </div>
    </section>
  );
}

// ---------- SECTION: FINAL CTA ----------

function FinalCta() {
  return (
    <section className="flex w-full flex-col items-center bg-[var(--background)] px-5 pt-[100px] lg:px-20 lg:pt-[140px]">
      <div className="flex w-full max-w-[1760px] flex-col items-center gap-8 rounded-[32px] bg-gradient-to-br from-[var(--primary-warm)] via-[#F97316] to-[#EA580C] px-8 py-16 text-center lg:px-16 lg:py-20">
        <h2 className="text-[44px] font-extrabold leading-[110%] text-[var(--primary-foreground)] lg:text-[56px]">
          Готовы построить маршрут?
        </h2>
        <p className="max-w-[640px] text-[20px] leading-[160%] text-[var(--primary-foreground)]/90">
          Бесплатно для всей семьи. 5 минут на маршрут. Доплата — только за ИИ-ассистента. 200+
          семей.
        </p>
        <div className="inline-flex items-center gap-2 px-4 py-2">
          <Target className="size-4 text-[var(--primary-foreground)]" strokeWidth={2} />
          <span className="text-[15px] font-semibold text-[var(--primary-foreground)]">
            Не проходите курс — достигайте цель
          </span>
        </div>
        <button
          type="button"
          className="inline-flex items-center gap-3 rounded-full bg-[var(--primary-foreground)] px-8 py-4 text-[17px] font-bold text-[var(--primary-warm)] transition-opacity hover:opacity-90"
        >
          Начать строить маршрут
          <ArrowIconCircle />
        </button>
      </div>
    </section>
  );
}

// ---------- SECTION: FOOTER ----------

const FOOTER_COLUMNS = [
  {
    title: 'Продукт',
    links: ['Как это работает', 'Возможности', 'Тарифы', 'Отзывы'],
  },
  {
    title: 'Поддержка',
    links: ['FAQ', 'Контакты', 'Блог'],
  },
  {
    title: 'Правовая информация',
    links: ['Политика конфиденциальности', 'Условия использования', 'Согласие на обработку данных'],
  },
  {
    title: 'Для бизнеса',
    links: ['Для бизнеса', 'VEDO Hub'],
  },
];

function Footer() {
  return (
    <footer className="flex w-full flex-col bg-[var(--card)] px-5 py-16 lg:px-20">
      <div className="mx-auto flex w-full max-w-[1760px] flex-col gap-12">
        <div className="flex flex-col justify-between gap-10 lg:flex-row lg:gap-0">
          <div className="flex flex-col gap-4">
            <div className="flex items-center gap-3">
              <IconHand />
              <span className="text-[22px] font-bold text-[var(--foreground)]">Дай пять</span>
            </div>
            <p className="max-w-[320px] text-[16px] leading-[160%] text-[var(--muted-foreground)]">
              Индивидуальные маршруты обучения на основе карты знаний для семей на семейном
              образовании.
            </p>
            <div className="flex items-center gap-3">
              <button
                type="button"
                aria-label="Send"
                className="flex size-10 items-center justify-center rounded-full bg-[var(--muted)] text-[var(--muted-foreground)] transition-colors hover:text-[var(--foreground)]"
              >
                <Send className="size-4" strokeWidth={2} />
              </button>
              <button
                type="button"
                aria-label="Users"
                className="flex size-10 items-center justify-center rounded-full bg-[var(--muted)] text-[var(--muted-foreground)] transition-colors hover:text-[var(--foreground)]"
              >
                <Users className="size-4" strokeWidth={2} />
              </button>
              <button
                type="button"
                aria-label="Play"
                className="flex size-10 items-center justify-center rounded-full bg-[var(--muted)] text-[var(--muted-foreground)] transition-colors hover:text-[var(--foreground)]"
              >
                <Play className="size-4" strokeWidth={2} />
              </button>
            </div>
          </div>
          <div className="grid grid-cols-2 gap-8 md:grid-cols-4 lg:gap-16">
            {FOOTER_COLUMNS.map((col) => (
              <div key={col.title} className="flex flex-col gap-4">
                <span className="text-[16px] font-bold text-[var(--foreground)]">{col.title}</span>
                <ul className="flex flex-col gap-3">
                  {col.links.map((link) => (
                    <li key={link}>
                      <a
                        href="/"
                        className="text-[15px] text-[var(--muted-foreground)] transition-colors hover:text-[var(--foreground)]"
                      >
                        {link}
                      </a>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </div>
        <div className="flex flex-col items-center justify-between gap-4 border-t border-[var(--border)] pt-8 sm:flex-row">
          <span className="text-[14px] text-[var(--muted-foreground)]">
            © 2026 «Дай пять» — сервис VEDO EduTrack
          </span>
          <span className="text-[14px] text-[var(--muted-foreground)]">
            Сделано на основе открытой онтологии VEDO Hub
          </span>
        </div>
      </div>
    </footer>
  );
}

// ---------- PAGE ----------

export function LandingPage() {
  useTheme();
  return (
    <div className="w-full bg-[var(--background)] text-[var(--foreground)]">
      <Header />
      <main className="mx-auto w-full max-w-[1920px] overflow-x-hidden">
        <Hero />
        <ProblemSection />
        <Bridge />
        <SolutionSection />
        <BenefitsSection />
        <PricingSection />
        <TestimonialsSection />
        <PhilosophySection />
        <FaqSection />
        <FinalCta />
      </main>
      <Footer />
      <button
        type="button"
        onClick={() => {
          const current = document.documentElement.getAttribute('data-collection-3-4-mode');
          const next = current === 'dark' ? 'light' : 'dark';
          document.documentElement.setAttribute('data-collection-3-4-mode', next);
          localStorage.setItem(THEME_KEY, next);
        }}
        className="fixed bottom-6 right-6 z-50 flex size-10 items-center justify-center rounded-full bg-[var(--card)] text-[var(--foreground)] shadow-[0_12px_48px_-16px_rgba(124,45,18,0.35)]"
        aria-label="Переключить тему"
      >
        {document.documentElement.getAttribute('data-collection-3-4-mode') === 'dark' ? (
          <Sparkles className="size-4" strokeWidth={2} />
        ) : (
          <HelpCircle className="size-4" strokeWidth={2} />
        )}
      </button>
    </div>
  );
}
