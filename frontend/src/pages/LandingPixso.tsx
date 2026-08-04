import { useEffect } from 'react';
import '../styles/pixso-variables.css';
import '../styles/pixso-landing.css';

const THEME_KEY = 'pixso-theme';

function getInitialTheme(): string {
  const stored = localStorage.getItem(THEME_KEY);
  if (stored === 'dark' || stored === 'light') return stored;
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function CTAButton(props: any) {
  const { children } = props;
  return (
    <div className="component-4_257">
      <div className="Pixso-symbol-4_257">
        {children ?? <p className="Pixso-paragraph-4_258">Начать бесплатно</p>}
        <div className="Pixso-frame-4_259">
          <div className="frame-content-4_259">
            <div className="Pixso-frame-4_260">
              <div className="Pixso-vector-4_261"></div>
              <div className="Pixso-vector-4_262"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function Eyebrow(props: any) {
  const { children } = props;
  return (
    <div className="component-4_271">
      <div className="Pixso-symbol-4_271">
        <div className="Pixso-frame-4_272"></div>
        {children ?? <p className="Pixso-paragraph-4_273">СЕКЦИЯ</p>}
      </div>
    </div>
  );
}

function SectionHeader(props: any) {
  const { eyebrow, title, subtitle } = props;
  return (
    <div className="component-7_207">
      <div className="Pixso-symbol-7_207">
        {eyebrow ?? <Eyebrow />}
        {title ?? <p className="Pixso-paragraph-7_211">Заголовок секции</p>}
        {subtitle ?? <p className="Pixso-paragraph-7_212">Подзаголовок секции</p>}
      </div>
    </div>
  );
}

export function LandingPixso() {
  useEffect(() => {
    const theme = getInitialTheme();
    document.documentElement.setAttribute('data-collection-3-4-mode', theme);

    const mq = window.matchMedia('(prefers-color-scheme: dark)');
    const handler = (e: MediaQueryListEvent) => {
      const next = e.matches ? 'dark' : 'light';
      document.documentElement.setAttribute('data-collection-3-4-mode', next);
      localStorage.setItem(THEME_KEY, next);
    };
    mq.addEventListener('change', handler);
    return () => mq.removeEventListener('change', handler);
  }, []);

  return (
    <div className="scroll-container">
          <div id="3_1309" className="Pixso-frame-3_1309">
                <div id="3_1310" className="Pixso-frame-3_1310">
                    <div className="frame-content-3_1310">
                        <div id="4_99" className="stroke-wrapper-4_99">
                            <div className="Pixso-frame-4_99"></div>
                            <div className="stroke-4_99"></div>
                            <div className="Pixso-frame-4_99-content-layer">
                                <div className="frame-content-4_99">
                                    <div
                                        id="3_1311"
                                        className="Pixso-frame-3_1311"
                                    >
                                        <div
                                            id="3_1312"
                                            className="Pixso-vector-3_1312"
                                        ></div>
                                        <p
                                            id="3_1317"
                                            className="Pixso-paragraph-3_1317"
                                        >
                                            {"Дай пять"}
                                        </p>
                                    </div>
                                    <div
                                        id="3_1318"
                                        className="Pixso-frame-3_1318"
                                    >
                                        <p
                                            id="3_1319"
                                            className="Pixso-paragraph-3_1319"
                                        >
                                            {"Как это работает"}
                                        </p>
                                        <p
                                            id="3_1320"
                                            className="Pixso-paragraph-3_1320"
                                        >
                                            {"Возможности"}
                                        </p>
                                        <p
                                            id="6_120"
                                            className="Pixso-paragraph-6_120"
                                        >
                                            {"Тарифы"}
                                        </p>
                                        <p
                                            id="3_1321"
                                            className="Pixso-paragraph-3_1321"
                                        >
                                            {"Отзывы"}
                                        </p>
                                        <p
                                            id="3_1322"
                                            className="Pixso-paragraph-3_1322"
                                        >
                                            {"FAQ"}
                                        </p>
                                    </div>
                                    <div id="7_7" className="Pixso-frame-7_7">
                                        <div
                                            id="7_8"
                                            className="Pixso-frame-7_8"
                                        >
                                            <div
                                                id="7_9"
                                                className="Pixso-vector-7_9"
                                            ></div>
                                            <p
                                                id="7_14"
                                                className="Pixso-paragraph-7_14"
                                            >
                                                {"Для бизнеса"}
                                            </p>
                                        </div>
                                        <div
                                            id="8_6"
                                            className="Pixso-frame-8_6"
                                        >
                                            <div
                                                id="8_7"
                                                className="Pixso-vector-8_7"
                                            ></div>
                                            <p
                                                id="8_11"
                                                className="Pixso-paragraph-8_11"
                                            >
                                                {"VEDO Hub"}
                                            </p>
                                        </div>
                                    </div>
                                    <CTAButton
                                        id="7_291"
                                        className="Pixso-instance-7_291"
                                        slot_4_258={
                                            <p
                                                id="11_38"
                                                className="Pixso-paragraph-11_38"
                                            >
                                                {"Попробовать бесплатно"}
                                            </p>
                                        }
                                    ></CTAButton>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="3_1325" className="Pixso-frame-3_1325">
                    <div className="frame-content-3_1325">
                        <div id="3_1326" className="Pixso-frame-3_1326">
                            <div className="frame-content-3_1326">
                                <div id="3_1327" className="Pixso-frame-3_1327">
                                    <div
                                        id="3_1328"
                                        className="Pixso-vector-3_1328"
                                    ></div>
                                    <p
                                        id="3_1333"
                                        className="Pixso-paragraph-3_1333"
                                    >
                                        {"Не зубрёжка, а понимание"}
                                    </p>
                                </div>
                                <div id="3_1334" className="Pixso-text-3_1334">
                                    <p
                                        id="3_1334_0"
                                        className="Pixso-paragraph-3_1334_0"
                                    >
                                        <span
                                            id="3_1334_0_1"
                                            className="Pixso-span-3_1334_0_1"
                                        >
                                            {"Ребёнок учится думать."}
                                        </span>
                                    </p>
                                    <p
                                        id="3_1334_1"
                                        className="Pixso-paragraph-3_1334_1"
                                    >
                                        <span
                                            id="3_1334_1_1"
                                            className="Pixso-span-3_1334_1_1"
                                        >
                                            {"Вы — направлять."}
                                        </span>
                                    </p>
                                </div>
                                <p
                                    id="3_1335"
                                    className="Pixso-paragraph-3_1335"
                                >
                                    {
                                        "Платформа для построения индивидуальных маршрутов обучения на основе графа знаний, диагностики пробелов, отслеживания прогресса и ИИ-ассистента"
                                    }
                                </p>
                                <div id="3_1336" className="Pixso-frame-3_1336">
                                    <CTAButton
                                        id="7_297"
                                        className="Pixso-instance-7_297"
                                        slot_4_258={
                                            <p
                                                id="11_33"
                                                className="Pixso-paragraph-11_33"
                                            >
                                                {"Построить маршрут за 5 минут"}
                                            </p>
                                        }
                                    ></CTAButton>
                                    <div
                                        id="3_1344"
                                        className="Pixso-frame-3_1344"
                                    >
                                        <p
                                            id="3_1345"
                                            className="Pixso-paragraph-3_1345"
                                        >
                                            {"Узнать больше"}
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div id="3_1346" className="Pixso-frame-3_1346">
                            <div id="4_45" className="Pixso-vector-4_45"></div>
                            <div id="4_46" className="Pixso-vector-4_46"></div>
                            <div id="4_101" className="stroke-wrapper-4_101">
                                <div className="Pixso-frame-4_101">
                                    <div className="shadow-blend-unknown-1"></div>
                                    <div className="shadow-blend-unknown-0"></div>
                                </div>
                                <div className="stroke-4_101"></div>
                                <div className="Pixso-frame-4_101-content-layer">
                                    <div className="frame-content-4_101">
                                        <div
                                            id="4_102"
                                            className="Pixso-frame-4_102"
                                        >
                                            <div className="frame-content-4_102">
                                                <p
                                                    id="4_103"
                                                    className="Pixso-paragraph-4_103"
                                                >
                                                    {"Карта знаний"}
                                                </p>
                                                <div
                                                    id="4_104"
                                                    className="Pixso-frame-4_104"
                                                >
                                                    <div
                                                        id="4_105"
                                                        className="Pixso-vector-4_105"
                                                    ></div>
                                                    <p
                                                        id="4_107"
                                                        className="Pixso-paragraph-4_107"
                                                    >
                                                        {"ФГОС ✓"}
                                                    </p>
                                                </div>
                                            </div>
                                        </div>
                                        <div
                                            id="4_108"
                                            className="Pixso-frame-4_108"
                                        >
                                            <div className="frame-content-4_108">
                                                <p
                                                    id="4_109"
                                                    className="Pixso-paragraph-4_109"
                                                >
                                                    {"Алгебра: деление → дроби"}
                                                </p>
                                                <div
                                                    id="4_110"
                                                    className="Pixso-frame-4_110"
                                                >
                                                    <div className="frame-content-4_110">
                                                        <div
                                                            id="4_111"
                                                            className="Pixso-frame-4_111"
                                                        ></div>
                                                    </div>
                                                </div>
                                                <div
                                                    id="4_112"
                                                    className="Pixso-frame-4_112"
                                                >
                                                    <div className="frame-content-4_112">
                                                        <p
                                                            id="4_113"
                                                            className="Pixso-paragraph-4_113"
                                                        >
                                                            {"Пройдено 75%"}
                                                        </p>
                                                        <p
                                                            id="4_114"
                                                            className="Pixso-paragraph-4_114"
                                                        >
                                                            {"12/16 тем"}
                                                        </p>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                        <div
                                            id="4_115"
                                            className="stroke-wrapper-4_115"
                                        >
                                            <div className="Pixso-frame-4_115">
                                                <div className="frame-content-4_115">
                                                    <div
                                                        id="4_116"
                                                        className="Pixso-vector-4_116"
                                                    ></div>
                                                    <div
                                                        id="4_120"
                                                        className="Pixso-frame-4_120"
                                                    >
                                                        <div className="frame-content-4_120">
                                                            <p
                                                                id="4_121"
                                                                className="Pixso-paragraph-4_121"
                                                            >
                                                                {
                                                                    "Пробел: Проценты"
                                                                }
                                                            </p>
                                                            <p
                                                                id="4_122"
                                                                className="Pixso-paragraph-4_122"
                                                            >
                                                                {
                                                                    "Нужны деление → дроби"
                                                                }
                                                            </p>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                            <div className="stroke-4_115"></div>
                                        </div>
                                        <div
                                            id="4_123"
                                            className="Pixso-frame-4_123"
                                        >
                                            <div className="frame-content-4_123">
                                                <p
                                                    id="4_124"
                                                    className="Pixso-paragraph-4_124"
                                                >
                                                    {"Маршрут до аттестации"}
                                                </p>
                                                <div
                                                    id="4_125"
                                                    className="Pixso-frame-4_125"
                                                >
                                                    <div className="frame-content-4_125">
                                                        <div
                                                            id="4_126"
                                                            className="Pixso-frame-4_126"
                                                        >
                                                            <div className="frame-content-4_126">
                                                                <p
                                                                    id="4_127"
                                                                    className="Pixso-paragraph-4_127"
                                                                >
                                                                    {"Деление"}
                                                                </p>
                                                            </div>
                                                        </div>
                                                        <div
                                                            id="4_128"
                                                            className="Pixso-vector-4_128"
                                                        ></div>
                                                        <div
                                                            id="4_131"
                                                            className="Pixso-frame-4_131"
                                                        >
                                                            <div className="frame-content-4_131">
                                                                <p
                                                                    id="4_132"
                                                                    className="Pixso-paragraph-4_132"
                                                                >
                                                                    {"Дроби"}
                                                                </p>
                                                            </div>
                                                        </div>
                                                        <div
                                                            id="4_133"
                                                            className="Pixso-vector-4_133"
                                                        ></div>
                                                        <div
                                                            id="4_136"
                                                            className="Pixso-frame-4_136"
                                                        >
                                                            <div className="frame-content-4_136">
                                                                <p
                                                                    id="4_137"
                                                                    className="Pixso-paragraph-4_137"
                                                                >
                                                                    {"Проценты"}
                                                                </p>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="3_1347" className="Pixso-frame-3_1347">
                    <div className="frame-content-3_1347">
                        <SectionHeader
                            id="7_213"
                            className="Pixso-instance-7_213"
                            slot_7_212={
                                <p
                                    id="11_126"
                                    className="Pixso-paragraph-11_126"
                                >
                                    {
                                        "Знакомые ситуации для родителей на семейном образовании"
                                    }
                                </p>
                            }
                            slot_7_211={
                                <p
                                    id="11_125"
                                    className="Pixso-paragraph-11_125"
                                >
                                    {"Это про вас?"}
                                </p>
                            }
                            slot_7_208={
                                <Eyebrow
                                    id="11_122"
                                    className="Pixso-instance-11_122"
                                    slot_4_273={
                                        <p
                                            id="11_124"
                                            className="Pixso-paragraph-11_124"
                                        >
                                            {"ПРОБЛЕМА"}
                                        </p>
                                    }
                                ></Eyebrow>
                            }
                        ></SectionHeader>
                        <div id="3_1351" className="Pixso-frame-3_1351">
                            <div className="frame-content-3_1351">
                                <div
                                    id="3_1352"
                                    className="stroke-wrapper-3_1352"
                                >
                                    <div className="Pixso-frame-3_1352">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1352">
                                            <div
                                                id="3_1353"
                                                className="Pixso-vector-3_1353"
                                            ></div>
                                            <p
                                                id="3_1354"
                                                className="Pixso-paragraph-3_1354"
                                            >
                                                {"Программа перегружена"}
                                            </p>
                                            <p
                                                id="3_1355"
                                                className="Pixso-paragraph-3_1355"
                                            >
                                                {
                                                    "Школьную программу невозможно осмысленно пройти в отведённое время: либо поверхностно, либо не успеваем."
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1352"></div>
                                </div>
                                <div
                                    id="3_1356"
                                    className="stroke-wrapper-3_1356"
                                >
                                    <div className="Pixso-frame-3_1356">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1356">
                                            <div
                                                id="3_1357"
                                                className="Pixso-vector-3_1357"
                                            ></div>
                                            <p
                                                id="3_1360"
                                                className="Pixso-paragraph-3_1360"
                                            >
                                                {"Пробелы в знаниях"}
                                            </p>
                                            <p
                                                id="3_1361"
                                                className="Pixso-paragraph-3_1361"
                                            >
                                                {
                                                    "В декабре выяснилось, что ребёнок не знает тему, которой не было в купленных курсах. Где ещё пробелы?"
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1356"></div>
                                </div>
                                <div
                                    id="3_1362"
                                    className="stroke-wrapper-3_1362"
                                >
                                    <div className="Pixso-frame-3_1362">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1362">
                                            <div
                                                id="3_1363"
                                                className="Pixso-vector-3_1363"
                                            ></div>
                                            <p
                                                id="3_1370"
                                                className="Pixso-paragraph-3_1370"
                                            >
                                                {"Предметы изолированы"}
                                            </p>
                                            <p
                                                id="3_1371"
                                                className="Pixso-paragraph-3_1371"
                                            >
                                                {
                                                    "Физика отдельно от математики, биология — от химии. Ребёнок не видит связей между явлениями."
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1362"></div>
                                </div>
                                <div
                                    id="3_1372"
                                    className="stroke-wrapper-3_1372"
                                >
                                    <div className="Pixso-frame-3_1372">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1372">
                                            <div
                                                id="3_1373"
                                                className="Pixso-vector-3_1373"
                                            ></div>
                                            <p
                                                id="3_1376"
                                                className="Pixso-paragraph-3_1376"
                                            >
                                                {"Единый темп не работает"}
                                            </p>
                                            <p
                                                id="3_1377"
                                                className="Pixso-paragraph-3_1377"
                                            >
                                                {
                                                    "По математике — 7 класс, по биологии — 8. Это нормально. Но как это учесть в плане?"
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1372"></div>
                                </div>
                                <div
                                    id="3_1378"
                                    className="stroke-wrapper-3_1378"
                                >
                                    <div className="Pixso-frame-3_1378">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1378">
                                            <div
                                                id="3_1379"
                                                className="Pixso-vector-3_1379"
                                            ></div>
                                            <p
                                                id="3_1380"
                                                className="Pixso-paragraph-3_1380"
                                            >
                                                {"Непонятно зачем"}
                                            </p>
                                            <p
                                                id="3_1381"
                                                className="Pixso-paragraph-3_1381"
                                            >
                                                {
                                                    "Ребёнок спрашивает: «Зачем мне это?» — а ответить нечего. Каждое знание должно иметь смысл."
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1378"></div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="3_1391" className="Pixso-frame-3_1391">
                    <div className="frame-content-3_1391">
                        <p id="3_1392" className="Pixso-paragraph-3_1392">
                            {"Так было."}
                        </p>
                        <p id="3_1393" className="Pixso-paragraph-3_1393">
                            {"А теперь — так"}
                        </p>
                    </div>
                </div>
                <div id="3_1437" className="Pixso-frame-3_1437">
                    <div className="frame-content-3_1437">
                        <SectionHeader
                            id="7_219"
                            className="Pixso-instance-7_219"
                            slot_7_212={
                                <p
                                    id="11_141"
                                    className="Pixso-paragraph-11_141"
                                >
                                    {
                                        "Представьте, что все школьные темы — это города на карте, а связи между ними — дороги. Мы ведём вас по маршруту."
                                    }
                                </p>
                            }
                            slot_7_211={
                                <p
                                    id="11_140"
                                    className="Pixso-paragraph-11_140"
                                >
                                    {"Как работает «Дай пять»"}
                                </p>
                            }
                            slot_7_208={
                                <Eyebrow
                                    id="11_137"
                                    className="Pixso-instance-11_137"
                                    slot_4_273={
                                        <p
                                            id="11_139"
                                            className="Pixso-paragraph-11_139"
                                        >
                                            {"РЕШЕНИЕ"}
                                        </p>
                                    }
                                ></Eyebrow>
                            }
                        ></SectionHeader>
                        <div id="3_1441" className="Pixso-frame-3_1441">
                            <div className="frame-content-3_1441">
                                <div
                                    id="3_1442"
                                    className="stroke-wrapper-3_1442"
                                >
                                    <div className="Pixso-frame-3_1442">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1442">
                                            <div
                                                id="3_1443"
                                                className="Pixso-frame-3_1443"
                                            >
                                                <div className="frame-content-3_1443">
                                                    <p
                                                        id="3_1444"
                                                        className="Pixso-paragraph-3_1444"
                                                    >
                                                        {"1"}
                                                    </p>
                                                </div>
                                            </div>
                                            <p
                                                id="3_1445"
                                                className="Pixso-paragraph-3_1445"
                                            >
                                                {"Задаём цели и ограничения"}
                                            </p>
                                            <p
                                                id="3_1446"
                                                className="Pixso-paragraph-3_1446"
                                            >
                                                {
                                                    "Цели к аттестации, время, нагрузка, темп. ИИ-ассистент поможет сформулировать и скорректировать."
                                                }
                                            </p>
                                            <p
                                                id="3_1447"
                                                className="Pixso-paragraph-3_1447"
                                            >
                                                {
                                                    "Пример: «К маю — аттестация, по алгебре подтянуть до «4»»"
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1442"></div>
                                </div>
                                <div
                                    id="3_1448"
                                    className="stroke-wrapper-3_1448"
                                >
                                    <div className="Pixso-frame-3_1448">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1448">
                                            <div
                                                id="3_1449"
                                                className="Pixso-frame-3_1449"
                                            >
                                                <div className="frame-content-3_1449">
                                                    <p
                                                        id="3_1450"
                                                        className="Pixso-paragraph-3_1450"
                                                    >
                                                        {"2"}
                                                    </p>
                                                </div>
                                            </div>
                                            <p
                                                id="3_1451"
                                                className="Pixso-paragraph-3_1451"
                                            >
                                                {"Подключаем карту знаний"}
                                            </p>
                                            <p
                                                id="3_1452"
                                                className="Pixso-paragraph-3_1452"
                                            >
                                                {
                                                    "Подключаем готовую карту от методистов и сообщества: темы, связи, актуальный ФГОС."
                                                }
                                            </p>
                                            <p
                                                id="3_1453"
                                                className="Pixso-paragraph-3_1453"
                                            >
                                                {
                                                    "Пример: «Проценты → биология → химия → география»"
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1448"></div>
                                </div>
                                <div
                                    id="3_1454"
                                    className="stroke-wrapper-3_1454"
                                >
                                    <div className="Pixso-frame-3_1454">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1454">
                                            <div
                                                id="3_1455"
                                                className="Pixso-frame-3_1455"
                                            >
                                                <div className="frame-content-3_1455">
                                                    <p
                                                        id="3_1456"
                                                        className="Pixso-paragraph-3_1456"
                                                    >
                                                        {"3"}
                                                    </p>
                                                </div>
                                            </div>
                                            <p
                                                id="3_1457"
                                                className="Pixso-paragraph-3_1457"
                                            >
                                                {"Находим пробелы"}
                                            </p>
                                            <p
                                                id="3_1458"
                                                className="Pixso-paragraph-3_1458"
                                            >
                                                {
                                                    "Система показывает, какие темы пропущены и почему это важно."
                                                }
                                            </p>
                                            <p
                                                id="3_1459"
                                                className="Pixso-paragraph-3_1459"
                                            >
                                                {
                                                    "Пример: «Деление → Дроби → Проценты. Без деления — пробел»"
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1454"></div>
                                </div>
                                <div
                                    id="3_1460"
                                    className="stroke-wrapper-3_1460"
                                >
                                    <div className="Pixso-frame-3_1460">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1460">
                                            <div
                                                id="3_1461"
                                                className="Pixso-frame-3_1461"
                                            >
                                                <div className="frame-content-3_1461">
                                                    <p
                                                        id="3_1462"
                                                        className="Pixso-paragraph-3_1462"
                                                    >
                                                        {"4"}
                                                    </p>
                                                </div>
                                            </div>
                                            <p
                                                id="3_1463"
                                                className="Pixso-paragraph-3_1463"
                                            >
                                                {"Видим прогноз до аттестации"}
                                            </p>
                                            <p
                                                id="3_1464"
                                                className="Pixso-paragraph-3_1464"
                                            >
                                                {
                                                    "Видите, что успеете до экзамена. Планируете с уверенностью."
                                                }
                                            </p>
                                            <p
                                                id="3_1465"
                                                className="Pixso-paragraph-3_1465"
                                            >
                                                {
                                                    "Пример: «До ОГЭ 8 месяцев. Успеете 95% программы»"
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1460"></div>
                                </div>
                            </div>
                        </div>
                        <div id="4_348" className="Pixso-frame-4_348">
                            <div className="frame-content-4_348">
                                <div id="4_349" className="Pixso-frame-4_349">
                                    <p
                                        id="4_350"
                                        className="Pixso-paragraph-4_350"
                                    >
                                        {"Не знаете, с чего начать?"}
                                    </p>
                                    <p
                                        id="4_479"
                                        className="Pixso-paragraph-4_479"
                                    >
                                        {"Выберите педагогическую концепцию"}
                                    </p>
                                </div>
                                <div id="4_353" className="Pixso-frame-4_353">
                                    <div className="frame-content-4_353">
                                        <div
                                            id="4_354"
                                            className="stroke-wrapper-4_354"
                                        >
                                            <div className="Pixso-frame-4_354">
                                                <div className="frame-content-4_354">
                                                    <div
                                                        id="4_355"
                                                        className="Pixso-vector-4_355"
                                                    ></div>
                                                    <p
                                                        id="4_360"
                                                        className="Pixso-paragraph-4_360"
                                                    >
                                                        {"Спиральное обучение"}
                                                    </p>
                                                    <p
                                                        id="4_361"
                                                        className="Pixso-paragraph-4_361"
                                                    >
                                                        {
                                                            "Тема «Функции» возвращается в 7, 8, 9 классах — с возрастающей глубиной. Маршрут строит витки, а не линейную цепочку."
                                                        }
                                                    </p>
                                                    <div
                                                        id="4_362"
                                                        className="Pixso-frame-4_362"
                                                    >
                                                        <div
                                                            id="4_363"
                                                            className="Pixso-vector-4_363"
                                                        ></div>
                                                        <p
                                                            id="4_365"
                                                            className="Pixso-paragraph-4_365"
                                                        >
                                                            {"Выбрано"}
                                                        </p>
                                                    </div>
                                                </div>
                                            </div>
                                            <div className="stroke-4_354"></div>
                                        </div>
                                        <div
                                            id="4_366"
                                            className="stroke-wrapper-4_366"
                                        >
                                            <div className="Pixso-frame-4_366">
                                                <div className="frame-content-4_366">
                                                    <div
                                                        id="4_367"
                                                        className="Pixso-vector-4_367"
                                                    ></div>
                                                    <p
                                                        id="4_372"
                                                        className="Pixso-paragraph-4_372"
                                                    >
                                                        {"Проектное погружение"}
                                                    </p>
                                                    <p
                                                        id="4_373"
                                                        className="Pixso-paragraph-4_373"
                                                    >
                                                        {
                                                            "Модули из разных предметов группируются вокруг проекта «Экосистема»: биология + математика + химия + география."
                                                        }
                                                    </p>
                                                    <div
                                                        id="4_374"
                                                        className="Pixso-frame-4_374"
                                                    >
                                                        <p
                                                            id="4_375"
                                                            className="Pixso-paragraph-4_375"
                                                        >
                                                            {"Выбрать"}
                                                        </p>
                                                    </div>
                                                </div>
                                            </div>
                                            <div className="stroke-4_366"></div>
                                        </div>
                                        <div
                                            id="4_376"
                                            className="stroke-wrapper-4_376"
                                        >
                                            <div className="Pixso-frame-4_376">
                                                <div className="frame-content-4_376">
                                                    <div
                                                        id="4_377"
                                                        className="Pixso-vector-4_377"
                                                    ></div>
                                                    <p
                                                        id="4_381"
                                                        className="Pixso-paragraph-4_381"
                                                    >
                                                        {"Свободное обучение"}
                                                    </p>
                                                    <p
                                                        id="4_382"
                                                        className="Pixso-paragraph-4_382"
                                                    >
                                                        {
                                                            "Ребёнок идёт от интереса: динозавры → биология → география → химия. Маршрут строится вокруг увлечений, а не наоборот."
                                                        }
                                                    </p>
                                                    <div
                                                        id="4_383"
                                                        className="Pixso-frame-4_383"
                                                    >
                                                        <p
                                                            id="4_384"
                                                            className="Pixso-paragraph-4_384"
                                                        >
                                                            {"Выбрать"}
                                                        </p>
                                                    </div>
                                                </div>
                                            </div>
                                            <div className="stroke-4_376"></div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="3_1466" className="Pixso-frame-3_1466">
                    <div className="frame-content-3_1466">
                        <SectionHeader
                            id="7_225"
                            className="Pixso-instance-7_225"
                            slot_7_212={
                                <p
                                    id="11_136"
                                    className="Pixso-paragraph-11_136"
                                >
                                    {
                                        "Вместо метода проб и ошибок — система, которая работает"
                                    }
                                </p>
                            }
                            slot_7_211={
                                <p
                                    id="11_135"
                                    className="Pixso-paragraph-11_135"
                                >
                                    {"Почему родители выбирают «Дай пять»"}
                                </p>
                            }
                            slot_7_208={
                                <Eyebrow
                                    id="11_132"
                                    className="Pixso-instance-11_132"
                                    slot_4_273={
                                        <p
                                            id="11_134"
                                            className="Pixso-paragraph-11_134"
                                        >
                                            {"ПРЕИМУЩЕСТВА"}
                                        </p>
                                    }
                                ></Eyebrow>
                            }
                        ></SectionHeader>
                        <div id="3_1470" className="Pixso-frame-3_1470">
                            <div className="frame-content-3_1470">
                                <div id="3_1471" className="Pixso-frame-3_1471">
                                    <div className="frame-content-3_1471">
                                        <div
                                            id="3_1472"
                                            className="stroke-wrapper-3_1472"
                                        >
                                            <div className="Pixso-frame-3_1472">
                                                <div className="shadow-blend-unknown-1"></div>
                                                <div className="shadow-blend-unknown-0"></div>
                                                <div className="frame-content-3_1472">
                                                    <div
                                                        id="3_1473"
                                                        className="Pixso-vector-3_1473"
                                                    ></div>
                                                    <p
                                                        id="3_1476"
                                                        className="Pixso-paragraph-3_1476"
                                                    >
                                                        {
                                                            "Вижу пробелы до аттестации"
                                                        }
                                                    </p>
                                                    <p
                                                        id="3_1477"
                                                        className="Pixso-paragraph-3_1477"
                                                    >
                                                        {
                                                            "Система заранее покажет, какие темы не пройдены. Не будет сюрпризов на экзамене."
                                                        }
                                                    </p>
                                                    <p
                                                        id="3_1478"
                                                        className="Pixso-paragraph-3_1478"
                                                    >
                                                        {
                                                            "Вместо: «Надеюсь, всё успеем»"
                                                        }
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="stroke-3_1472"></div>
                                        </div>
                                        <div
                                            id="3_1479"
                                            className="stroke-wrapper-3_1479"
                                        >
                                            <div className="Pixso-frame-3_1479">
                                                <div className="shadow-blend-unknown-1"></div>
                                                <div className="shadow-blend-unknown-0"></div>
                                                <div className="frame-content-3_1479">
                                                    <div
                                                        id="3_1480"
                                                        className="Pixso-vector-3_1480"
                                                    ></div>
                                                    <p
                                                        id="3_1483"
                                                        className="Pixso-paragraph-3_1483"
                                                    >
                                                        {"Учу в своём темпе"}
                                                    </p>
                                                    <p
                                                        id="3_1484"
                                                        className="Pixso-paragraph-3_1484"
                                                    >
                                                        {
                                                            "По математике — 7 класс, по биологии — 8. Это нормально. Система учитывает разные уровни."
                                                        }
                                                    </p>
                                                    <p
                                                        id="3_1485"
                                                        className="Pixso-paragraph-3_1485"
                                                    >
                                                        {
                                                            "Вместо: «Все должны быть на одном уровне»"
                                                        }
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="stroke-3_1479"></div>
                                        </div>
                                        <div
                                            id="3_1486"
                                            className="stroke-wrapper-3_1486"
                                        >
                                            <div className="Pixso-frame-3_1486">
                                                <div className="shadow-blend-unknown-1"></div>
                                                <div className="shadow-blend-unknown-0"></div>
                                                <div className="frame-content-3_1486">
                                                    <div
                                                        id="3_1487"
                                                        className="Pixso-vector-3_1487"
                                                    ></div>
                                                    <p
                                                        id="3_1492"
                                                        className="Pixso-paragraph-3_1492"
                                                    >
                                                        {
                                                            "Одна панель для всех детей"
                                                        }
                                                    </p>
                                                    <p
                                                        id="3_1493"
                                                        className="Pixso-paragraph-3_1493"
                                                    >
                                                        {
                                                            "Трое детей? Не проблема. Видите прогресс каждого в одном месте."
                                                        }
                                                    </p>
                                                    <p
                                                        id="3_1494"
                                                        className="Pixso-paragraph-3_1494"
                                                    >
                                                        {
                                                            "Вместо: «Три таблицы Excel и два чата»"
                                                        }
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="stroke-3_1486"></div>
                                        </div>
                                    </div>
                                </div>
                                <div id="3_1495" className="Pixso-frame-3_1495">
                                    <div className="frame-content-3_1495">
                                        <div
                                            id="3_1496"
                                            className="stroke-wrapper-3_1496"
                                        >
                                            <div className="Pixso-frame-3_1496">
                                                <div className="shadow-blend-unknown-1"></div>
                                                <div className="shadow-blend-unknown-0"></div>
                                                <div className="frame-content-3_1496">
                                                    <div
                                                        id="3_1497"
                                                        className="Pixso-vector-3_1497"
                                                    ></div>
                                                    <p
                                                        id="3_1500"
                                                        className="Pixso-paragraph-3_1500"
                                                    >
                                                        {
                                                            "Знаю, зачем это учить"
                                                        }
                                                    </p>
                                                    <p
                                                        id="3_1501"
                                                        className="Pixso-paragraph-3_1501"
                                                    >
                                                        {
                                                            "Каждая тема связана с реальной жизнью и другими предметами. Ребёнок видит смысл."
                                                        }
                                                    </p>
                                                    <p
                                                        id="3_1502"
                                                        className="Pixso-paragraph-3_1502"
                                                    >
                                                        {
                                                            "Вместо: «Учи, потому что так надо»"
                                                        }
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="stroke-3_1496"></div>
                                        </div>
                                        <div
                                            id="3_1503"
                                            className="stroke-wrapper-3_1503"
                                        >
                                            <div className="Pixso-frame-3_1503">
                                                <div className="shadow-blend-unknown-1"></div>
                                                <div className="shadow-blend-unknown-0"></div>
                                                <div className="frame-content-3_1503">
                                                    <div
                                                        id="3_1504"
                                                        className="Pixso-vector-3_1504"
                                                    ></div>
                                                    <p
                                                        id="3_1508"
                                                        className="Pixso-paragraph-3_1508"
                                                    >
                                                        {
                                                            "Карта знаний показывает связи"
                                                        }
                                                    </p>
                                                    <p
                                                        id="3_1509"
                                                        className="Pixso-paragraph-3_1509"
                                                    >
                                                        {
                                                            "Визуальная карта с темами и связями между предметами. Физика не отдельно от математики."
                                                        }
                                                    </p>
                                                    <p
                                                        id="3_1510"
                                                        className="Pixso-paragraph-3_1510"
                                                    >
                                                        {
                                                            "Вместо: «Предметы как отдельные острова»"
                                                        }
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="stroke-3_1503"></div>
                                        </div>
                                        <div
                                            id="3_1511"
                                            className="stroke-wrapper-3_1511"
                                        >
                                            <div className="Pixso-frame-3_1511">
                                                <div className="shadow-blend-unknown-1"></div>
                                                <div className="shadow-blend-unknown-0"></div>
                                                <div className="frame-content-3_1511">
                                                    <div
                                                        id="3_1512"
                                                        className="Pixso-vector-3_1512"
                                                    ></div>
                                                    <p
                                                        id="3_1515"
                                                        className="Pixso-paragraph-3_1515"
                                                    >
                                                        {
                                                            "Успеваем подготовиться к аттестации"
                                                        }
                                                    </p>
                                                    <p
                                                        id="3_1516"
                                                        className="Pixso-paragraph-3_1516"
                                                    >
                                                        {
                                                            "Система показывает покрытие стандарта: какие темы закрыты, какие остались до аттестации. Без ручной сверки в Excel."
                                                        }
                                                    </p>
                                                    <p
                                                        id="3_1517"
                                                        className="Pixso-paragraph-3_1517"
                                                    >
                                                        {
                                                            "Вместо: «Ручная сверка с ФГОС в Excel»"
                                                        }
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="stroke-3_1511"></div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div id="3_1518" className="Pixso-frame-3_1518">
                            <div className="frame-content-3_1518">
                                <p
                                    id="3_1523"
                                    className="Pixso-paragraph-3_1523"
                                >
                                    {
                                        "Почему это работает: карта знаний основана на принципе межпредметных связей — каждая тема имеет понятные «зачем» и «где пригодится». А ИИ-ассистент ответит на вопросы и подскажет, как лучше выстроить маршрут."
                                    }
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="6_1" className="Pixso-frame-6_1">
                    <div className="frame-content-6_1">
                        <SectionHeader
                            id="7_249"
                            className="Pixso-instance-7_249"
                            slot_7_212={
                                <p
                                    id="11_111"
                                    className="Pixso-paragraph-11_111"
                                >
                                    {
                                        "Диагностика бесплатно. Мотивация — в Компасе. Экономия на репетиторах — в Навигаторе."
                                    }
                                </p>
                            }
                            slot_7_211={
                                <p
                                    id="11_110"
                                    className="Pixso-paragraph-11_110"
                                >
                                    {"Выберите свой маршрут к аттестации"}
                                </p>
                            }
                            slot_7_208={
                                <Eyebrow
                                    id="11_107"
                                    className="Pixso-instance-11_107"
                                    slot_4_273={
                                        <p
                                            id="11_109"
                                            className="Pixso-paragraph-11_109"
                                        >
                                            {"ТАРИФЫ"}
                                        </p>
                                    }
                                ></Eyebrow>
                            }
                        ></SectionHeader>
                        <div id="6_8" className="Pixso-frame-6_8">
                            <div className="frame-content-6_8">
                                <div id="7_23" className="Pixso-frame-7_23">
                                    <div className="frame-content-7_23">
                                        <div
                                            id="8_24"
                                            className="Pixso-frame-8_24"
                                        >
                                            <div
                                                id="8_25"
                                                className="Pixso-vector-8_25"
                                            ></div>
                                            <p
                                                id="7_24"
                                                className="Pixso-paragraph-7_24"
                                            >
                                                {"Карта"}
                                            </p>
                                        </div>
                                        <div
                                            id="7_25"
                                            className="Pixso-frame-7_25"
                                        >
                                            <p
                                                id="7_26"
                                                className="Pixso-paragraph-7_26"
                                            >
                                                {"0 ₽"}
                                            </p>
                                            <p
                                                id="7_27"
                                                className="Pixso-paragraph-7_27"
                                            >
                                                {"/ 7 дней"}
                                            </p>
                                        </div>
                                        <p
                                            id="7_28"
                                            className="Pixso-paragraph-7_28"
                                        >
                                            {
                                                "Диагностика пробелов и карта тем на 7 дней."
                                            }
                                        </p>
                                        <div
                                            id="7_29"
                                            className="Pixso-frame-7_29"
                                        ></div>
                                        <div
                                            id="7_30"
                                            className="Pixso-frame-7_30"
                                        >
                                            <div
                                                id="7_31"
                                                className="Pixso-frame-7_31"
                                            >
                                                <div
                                                    id="7_32"
                                                    className="Pixso-vector-7_32"
                                                ></div>
                                                <p
                                                    id="7_34"
                                                    className="Pixso-paragraph-7_34"
                                                >
                                                    {
                                                        "Карта всех тем по предметам"
                                                    }
                                                </p>
                                            </div>
                                            <div
                                                id="7_35"
                                                className="Pixso-frame-7_35"
                                            >
                                                <div
                                                    id="7_36"
                                                    className="Pixso-vector-7_36"
                                                ></div>
                                                <p
                                                    id="7_38"
                                                    className="Pixso-paragraph-7_38"
                                                >
                                                    {
                                                        "Диагностика пробелов до аттестации"
                                                    }
                                                </p>
                                            </div>
                                            <div
                                                id="7_626"
                                                className="Pixso-frame-7_626"
                                            >
                                                <div
                                                    id="7_627"
                                                    className="Pixso-vector-7_627"
                                                ></div>
                                                <p
                                                    id="7_629"
                                                    className="Pixso-paragraph-7_629"
                                                >
                                                    {"1 ребёнок в профиле"}
                                                </p>
                                            </div>
                                            <div
                                                id="7_43"
                                                className="Pixso-frame-7_43"
                                            >
                                                <div
                                                    id="7_44"
                                                    className="Pixso-vector-7_44"
                                                ></div>
                                                <p
                                                    id="7_46"
                                                    className="Pixso-paragraph-7_46"
                                                >
                                                    {
                                                        "Без привязки банковской карты"
                                                    }
                                                </p>
                                            </div>
                                            <div
                                                id="7_344"
                                                className="Pixso-frame-7_344"
                                            >
                                                <div
                                                    id="7_345"
                                                    className="Pixso-vector-7_345"
                                                ></div>
                                                <p
                                                    id="7_347"
                                                    className="Pixso-paragraph-7_347"
                                                >
                                                    {
                                                        "Истории «зачем это знать»"
                                                    }
                                                </p>
                                            </div>
                                            <div
                                                id="7_348"
                                                className="Pixso-frame-7_348"
                                            >
                                                <div
                                                    id="7_349"
                                                    className="Pixso-vector-7_349"
                                                ></div>
                                                <p
                                                    id="7_351"
                                                    className="Pixso-paragraph-7_351"
                                                >
                                                    {"Маршрут и расписание"}
                                                </p>
                                            </div>
                                            <div
                                                id="7_352"
                                                className="Pixso-frame-7_352"
                                            >
                                                <div
                                                    id="7_353"
                                                    className="Pixso-vector-7_353"
                                                ></div>
                                                <p
                                                    id="7_355"
                                                    className="Pixso-paragraph-7_355"
                                                >
                                                    {
                                                        "Планирование бюджета и ресурсов"
                                                    }
                                                </p>
                                            </div>
                                        </div>
                                        <div
                                            id="7_59"
                                            className="Pixso-frame-7_59"
                                        >
                                            <CTAButton
                                                id="7_60"
                                                className="Pixso-instance-7_60"
                                                slot_4_258={
                                                    <p
                                                        id="11_11"
                                                        className="Pixso-paragraph-11_11"
                                                    >
                                                        {"Попробовать 7 дней"}
                                                    </p>
                                                }
                                            ></CTAButton>
                                        </div>
                                    </div>
                                </div>
                                <div id="7_66" className="Pixso-frame-7_66">
                                    <div className="frame-content-7_66">
                                        <div
                                            id="8_29"
                                            className="Pixso-frame-8_29"
                                        >
                                            <div
                                                id="8_30"
                                                className="Pixso-vector-8_30"
                                            ></div>
                                            <p
                                                id="7_67"
                                                className="Pixso-paragraph-7_67"
                                            >
                                                {"Компас"}
                                            </p>
                                        </div>
                                        <div
                                            id="7_68"
                                            className="Pixso-frame-7_68"
                                        >
                                            <p
                                                id="7_69"
                                                className="Pixso-paragraph-7_69"
                                            >
                                                {"990 ₽"}
                                            </p>
                                            <p
                                                id="7_70"
                                                className="Pixso-paragraph-7_70"
                                            >
                                                {"/ месяц"}
                                            </p>
                                        </div>
                                        <div
                                            id="7_71"
                                            className="Pixso-frame-7_71"
                                        >
                                            <p
                                                id="7_72"
                                                className="Pixso-paragraph-7_72"
                                            >
                                                {"8 990 ₽"}
                                            </p>
                                            <p
                                                id="7_73"
                                                className="Pixso-paragraph-7_73"
                                            >
                                                {"11 880 ₽"}
                                            </p>
                                            <div
                                                id="7_74"
                                                className="Pixso-frame-7_74"
                                            >
                                                <p
                                                    id="7_75"
                                                    className="Pixso-paragraph-7_75"
                                                >
                                                    {"Экономия 2 890 ₽"}
                                                </p>
                                            </div>
                                        </div>
                                        <p
                                            id="7_76"
                                            className="Pixso-paragraph-7_76"
                                        >
                                            {
                                                "Направление для всей семьи. Истории, маршрут и расписание, 30 запросов к ИИ."
                                            }
                                        </p>
                                        <div
                                            id="7_77"
                                            className="Pixso-frame-7_77"
                                        ></div>
                                        <div
                                            id="7_78"
                                            className="Pixso-frame-7_78"
                                        >
                                            <div
                                                id="7_79"
                                                className="Pixso-frame-7_79"
                                            >
                                                <div
                                                    id="7_80"
                                                    className="Pixso-vector-7_80"
                                                ></div>
                                                <p
                                                    id="7_82"
                                                    className="Pixso-paragraph-7_82"
                                                >
                                                    {"Всё из Карты"}
                                                </p>
                                            </div>
                                            <div
                                                id="7_83"
                                                className="Pixso-frame-7_83"
                                            >
                                                <div
                                                    id="7_84"
                                                    className="Pixso-vector-7_84"
                                                ></div>
                                                <p
                                                    id="7_86"
                                                    className="Pixso-paragraph-7_86"
                                                >
                                                    {
                                                        "Истории и контекст «зачем это знать»"
                                                    }
                                                </p>
                                            </div>
                                            <div
                                                id="7_87"
                                                className="Pixso-frame-7_87"
                                            >
                                                <div
                                                    id="7_88"
                                                    className="Pixso-vector-7_88"
                                                ></div>
                                                <p
                                                    id="7_90"
                                                    className="Pixso-paragraph-7_90"
                                                >
                                                    {
                                                        "Маршрут и расписание: что доступно сейчас, что дальше, где цель"
                                                    }
                                                </p>
                                            </div>
                                            <div
                                                id="7_91"
                                                className="Pixso-frame-7_91"
                                            >
                                                <div
                                                    id="7_92"
                                                    className="Pixso-vector-7_92"
                                                ></div>
                                                <p
                                                    id="7_94"
                                                    className="Pixso-paragraph-7_94"
                                                >
                                                    {"Полный учебный маршрут"}
                                                </p>
                                            </div>
                                            <div
                                                id="7_95"
                                                className="Pixso-frame-7_95"
                                            >
                                                <div
                                                    id="7_96"
                                                    className="Pixso-vector-7_96"
                                                ></div>
                                                <p
                                                    id="7_98"
                                                    className="Pixso-paragraph-7_98"
                                                >
                                                    {"Отслеживание прогресса"}
                                                </p>
                                            </div>
                                            <div
                                                id="7_99"
                                                className="Pixso-frame-7_99"
                                            >
                                                <div
                                                    id="7_100"
                                                    className="Pixso-vector-7_100"
                                                ></div>
                                                <p
                                                    id="7_102"
                                                    className="Pixso-paragraph-7_102"
                                                >
                                                    {
                                                        "До 2 детей в одном профиле"
                                                    }
                                                </p>
                                            </div>
                                            <div
                                                id="7_103"
                                                className="Pixso-frame-7_103"
                                            >
                                                <div
                                                    id="7_104"
                                                    className="Pixso-vector-7_104"
                                                ></div>
                                                <p
                                                    id="7_106"
                                                    className="Pixso-paragraph-7_106"
                                                >
                                                    {
                                                        "ИИ-ассистент (30 запросов/мес)"
                                                    }
                                                </p>
                                            </div>
                                        </div>
                                        <div
                                            id="7_107"
                                            className="Pixso-frame-7_107"
                                        >
                                            <CTAButton
                                                id="7_108"
                                                className="Pixso-instance-7_108"
                                                slot_4_258={
                                                    <p
                                                        id="11_43"
                                                        className="Pixso-paragraph-11_43"
                                                    >
                                                        {"Выбрать маршрут"}
                                                    </p>
                                                }
                                            ></CTAButton>
                                        </div>
                                    </div>
                                </div>
                                <div id="7_114" className="Pixso-frame-7_114">
                                    <div className="frame-content-7_114">
                                        <div
                                            id="8_22"
                                            className="Pixso-frame-8_22"
                                        >
                                            <div
                                                id="8_33"
                                                className="Pixso-vector-8_33"
                                            ></div>
                                            <p
                                                id="7_117"
                                                className="Pixso-paragraph-7_117"
                                            >
                                                {"Навигатор"}
                                            </p>
                                            <div
                                                id="7_115"
                                                className="Pixso-frame-7_115"
                                            >
                                                <p
                                                    id="7_116"
                                                    className="Pixso-paragraph-7_116"
                                                >
                                                    {"Самый популярный"}
                                                </p>
                                            </div>
                                        </div>
                                        <div
                                            id="7_118"
                                            className="Pixso-frame-7_118"
                                        >
                                            <p
                                                id="7_119"
                                                className="Pixso-paragraph-7_119"
                                            >
                                                {"1 990 ₽"}
                                            </p>
                                            <p
                                                id="7_120"
                                                className="Pixso-paragraph-7_120"
                                            >
                                                {"/ месяц"}
                                            </p>
                                        </div>
                                        <div
                                            id="7_121"
                                            className="Pixso-frame-7_121"
                                        >
                                            <p
                                                id="7_122"
                                                className="Pixso-paragraph-7_122"
                                            >
                                                {"17 910 ₽"}
                                            </p>
                                            <p
                                                id="7_123"
                                                className="Pixso-paragraph-7_123"
                                            >
                                                {"23 880 ₽"}
                                            </p>
                                            <div
                                                id="7_124"
                                                className="Pixso-frame-7_124"
                                            >
                                                <p
                                                    id="7_125"
                                                    className="Pixso-paragraph-7_125"
                                                >
                                                    {"Экономия 5 970 ₽"}
                                                </p>
                                            </div>
                                        </div>
                                        <p
                                            id="7_126"
                                            className="Pixso-paragraph-7_126"
                                        >
                                            {
                                                "Полная система для осмысленной подготовки. 300 запросов к ИИ, ресурсы, проекты и отчёты."
                                            }
                                        </p>
                                        <div
                                            id="7_127"
                                            className="Pixso-frame-7_127"
                                        ></div>
                                        <div
                                            id="7_128"
                                            className="Pixso-frame-7_128"
                                        >
                                            <div
                                                id="7_129"
                                                className="Pixso-frame-7_129"
                                            >
                                                <div
                                                    id="7_130"
                                                    className="Pixso-vector-7_130"
                                                ></div>
                                                <p
                                                    id="7_132"
                                                    className="Pixso-paragraph-7_132"
                                                >
                                                    {"Всё из Компаса"}
                                                </p>
                                            </div>
                                            <div
                                                id="7_133"
                                                className="Pixso-frame-7_133"
                                            >
                                                <div
                                                    id="7_134"
                                                    className="Pixso-vector-7_134"
                                                ></div>
                                                <p
                                                    id="7_136"
                                                    className="Pixso-paragraph-7_136"
                                                >
                                                    {
                                                        "Без ограничений по количеству детей"
                                                    }
                                                </p>
                                            </div>
                                            <div
                                                id="7_137"
                                                className="Pixso-frame-7_137"
                                            >
                                                <div
                                                    id="7_138"
                                                    className="Pixso-vector-7_138"
                                                ></div>
                                                <p
                                                    id="7_140"
                                                    className="Pixso-paragraph-7_140"
                                                >
                                                    {
                                                        "ИИ-ассистент по планированию (300 запросов/мес)"
                                                    }
                                                </p>
                                            </div>
                                            <div
                                                id="7_141"
                                                className="Pixso-frame-7_141"
                                            >
                                                <div
                                                    id="7_142"
                                                    className="Pixso-vector-7_142"
                                                ></div>
                                                <p
                                                    id="7_144"
                                                    className="Pixso-paragraph-7_144"
                                                >
                                                    {
                                                        "Планирование бюджета и ресурсов — сколько денег и времени до аттестации"
                                                    }
                                                </p>
                                            </div>
                                            <div
                                                id="7_145"
                                                className="Pixso-frame-7_145"
                                            >
                                                <div
                                                    id="7_146"
                                                    className="Pixso-vector-7_146"
                                                ></div>
                                                <p
                                                    id="7_148"
                                                    className="Pixso-paragraph-7_148"
                                                >
                                                    {
                                                        "Проектные идеи на стыках предметов"
                                                    }
                                                </p>
                                            </div>
                                            <div
                                                id="7_149"
                                                className="Pixso-frame-7_149"
                                            >
                                                <div
                                                    id="7_150"
                                                    className="Pixso-vector-7_150"
                                                ></div>
                                                <p
                                                    id="7_152"
                                                    className="Pixso-paragraph-7_152"
                                                >
                                                    {
                                                        "Экспорт отчётов и сертификатов"
                                                    }
                                                </p>
                                            </div>
                                            <div
                                                id="7_153"
                                                className="Pixso-frame-7_153"
                                            >
                                                <div
                                                    id="7_154"
                                                    className="Pixso-vector-7_154"
                                                ></div>
                                                <p
                                                    id="7_156"
                                                    className="Pixso-paragraph-7_156"
                                                >
                                                    {"Приоритетная поддержка"}
                                                </p>
                                            </div>
                                        </div>
                                        <div
                                            id="7_157"
                                            className="Pixso-frame-7_157"
                                        >
                                            <CTAButton
                                                id="7_158"
                                                className="Pixso-instance-7_158"
                                                slot_4_258={
                                                    <p
                                                        id="11_18"
                                                        className="Pixso-paragraph-11_18"
                                                    >
                                                        {"Получить максимум"}
                                                    </p>
                                                }
                                            ></CTAButton>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div id="7_194" className="Pixso-frame-7_194">
                            <div className="frame-content-7_194">
                                <div
                                    id="7_195"
                                    className="Pixso-vector-7_195"
                                ></div>
                                <p id="7_196" className="Pixso-paragraph-7_196">
                                    {
                                        "Не нашли подходящий вариант? VEDO EduTrack также используют школы, онлайн-платформы и корпорации."
                                    }
                                </p>
                                <p id="7_197" className="Pixso-paragraph-7_197">
                                    {"Перейти в раздел «Для бизнеса»"}
                                </p>
                                <div
                                    id="7_198"
                                    className="Pixso-vector-7_198"
                                ></div>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="3_1524" className="Pixso-frame-3_1524">
                    <div className="frame-content-3_1524">
                        <SectionHeader
                            id="7_231"
                            className="Pixso-instance-7_231"
                            slot_7_212={
                                <p id="11_91" className="Pixso-paragraph-11_91">
                                    {
                                        "Как «Дай пять» помог семьям на семейном образовании"
                                    }
                                </p>
                            }
                            slot_7_211={
                                <p id="11_90" className="Pixso-paragraph-11_90">
                                    {"Реальные истории"}
                                </p>
                            }
                            slot_7_208={
                                <Eyebrow
                                    id="11_87"
                                    className="Pixso-instance-11_87"
                                    slot_4_273={
                                        <p
                                            id="11_89"
                                            className="Pixso-paragraph-11_89"
                                        >
                                            {"ОТЗЫВЫ"}
                                        </p>
                                    }
                                ></Eyebrow>
                            }
                        ></SectionHeader>
                        <div id="3_1528" className="Pixso-frame-3_1528">
                            <div className="frame-content-3_1528">
                                <div
                                    id="3_1529"
                                    className="stroke-wrapper-3_1529"
                                >
                                    <div className="Pixso-frame-3_1529">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1529">
                                            <div
                                                id="3_1530"
                                                className="Pixso-vector-3_1530"
                                            ></div>
                                            <p
                                                id="3_1533"
                                                className="Pixso-paragraph-3_1533"
                                            >
                                                {
                                                    "«До «Дай пять» я тратила 3 часа в неделю на составление плана. Сейчас система сама показывает, что учить дальше. Маша стала учиться с удовольствием — видит, зачем ей каждая тема.»"
                                                }
                                            </p>
                                            <div
                                                id="3_1534"
                                                className="Pixso-frame-3_1534"
                                            >
                                                <div className="frame-content-3_1534">
                                                    <div
                                                        id="3_1535"
                                                        className="Pixso-frame-3_1535"
                                                    ></div>
                                                    <div
                                                        id="3_1536"
                                                        className="Pixso-frame-3_1536"
                                                    >
                                                        <p
                                                            id="3_1537"
                                                            className="Pixso-paragraph-3_1537"
                                                        >
                                                            {"Елена, 38 лет"}
                                                        </p>
                                                        <p
                                                            id="3_1538"
                                                            className="Pixso-paragraph-3_1538"
                                                        >
                                                            {
                                                                "Мама Маши (6 класс) и Пети (3 класс)"
                                                            }
                                                        </p>
                                                    </div>
                                                </div>
                                            </div>
                                            <div
                                                id="3_1539"
                                                className="Pixso-frame-3_1539"
                                            >
                                                <div className="frame-content-3_1539">
                                                    <div
                                                        id="3_1540"
                                                        className="Pixso-vector-3_1540"
                                                    ></div>
                                                    <p
                                                        id="3_1542"
                                                        className="Pixso-paragraph-3_1542"
                                                    >
                                                        {
                                                            "Результат: аттестация сдана на «отлично», пробелов нет"
                                                        }
                                                    </p>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1529"></div>
                                </div>
                                <div
                                    id="3_1543"
                                    className="stroke-wrapper-3_1543"
                                >
                                    <div className="Pixso-frame-3_1543">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1543">
                                            <div
                                                id="3_1544"
                                                className="Pixso-vector-3_1544"
                                            ></div>
                                            <p
                                                id="3_1547"
                                                className="Pixso-paragraph-3_1547"
                                            >
                                                {
                                                    "«Сын обожает биологию, но алгебра шла туго. «Дай пять» показал: он не понимает дроби → проценты → статистику. Закрыли пробел за 2 месяца. Теперь сам просит задачи по математике!»"
                                                }
                                            </p>
                                            <div
                                                id="3_1548"
                                                className="Pixso-frame-3_1548"
                                            >
                                                <div className="frame-content-3_1548">
                                                    <div
                                                        id="3_1549"
                                                        className="Pixso-frame-3_1549"
                                                    ></div>
                                                    <div
                                                        id="3_1550"
                                                        className="Pixso-frame-3_1550"
                                                    >
                                                        <p
                                                            id="3_1551"
                                                            className="Pixso-paragraph-3_1551"
                                                        >
                                                            {"Андрей, 42 года"}
                                                        </p>
                                                        <p
                                                            id="3_1552"
                                                            className="Pixso-paragraph-3_1552"
                                                        >
                                                            {
                                                                "Папа Димы (8 класс)"
                                                            }
                                                        </p>
                                                    </div>
                                                </div>
                                            </div>
                                            <div
                                                id="3_1553"
                                                className="Pixso-frame-3_1553"
                                            >
                                                <div className="frame-content-3_1553">
                                                    <div
                                                        id="3_1554"
                                                        className="Pixso-vector-3_1554"
                                                    ></div>
                                                    <p
                                                        id="3_1556"
                                                        className="Pixso-paragraph-3_1556"
                                                    >
                                                        {
                                                            "Результат: оценка по алгебре выросла с 3 до 5"
                                                        }
                                                    </p>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1543"></div>
                                </div>
                                <div
                                    id="3_1557"
                                    className="stroke-wrapper-3_1557"
                                >
                                    <div className="Pixso-frame-3_1557">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1557">
                                            <div
                                                id="3_1558"
                                                className="Pixso-vector-3_1558"
                                            ></div>
                                            <p
                                                id="3_1561"
                                                className="Pixso-paragraph-3_1561"
                                            >
                                                {
                                                    "«Трое детей на семейном — это хаос. «Дай пять» собрал всё в одну панель: вижу прогресс каждого, дети сами планируют обучение.»"
                                                }
                                            </p>
                                            <div
                                                id="3_1562"
                                                className="Pixso-frame-3_1562"
                                            >
                                                <div className="frame-content-3_1562">
                                                    <div
                                                        id="3_1563"
                                                        className="Pixso-frame-3_1563"
                                                    ></div>
                                                    <div
                                                        id="3_1564"
                                                        className="Pixso-frame-3_1564"
                                                    >
                                                        <p
                                                            id="3_1565"
                                                            className="Pixso-paragraph-3_1565"
                                                        >
                                                            {"Ольга, 35 лет"}
                                                        </p>
                                                        <p
                                                            id="3_1566"
                                                            className="Pixso-paragraph-3_1566"
                                                        >
                                                            {
                                                                "Мама троих: Кати (9 класс), Вани (6 класс), Сони (2 класс)"
                                                            }
                                                        </p>
                                                    </div>
                                                </div>
                                            </div>
                                            <div
                                                id="3_1567"
                                                className="Pixso-frame-3_1567"
                                            >
                                                <div className="frame-content-3_1567">
                                                    <div
                                                        id="3_1568"
                                                        className="Pixso-vector-3_1568"
                                                    ></div>
                                                    <p
                                                        id="3_1570"
                                                        className="Pixso-paragraph-3_1570"
                                                    >
                                                        {
                                                            "Результат: дети учатся самостоятельно, мама — координатор"
                                                        }
                                                    </p>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1557"></div>
                                </div>
                            </div>
                        </div>
                        <div id="3_1571" className="Pixso-frame-3_1571">
                            <div className="frame-content-3_1571">
                                <div id="3_1572" className="Pixso-frame-3_1572">
                                    <p
                                        id="3_1573"
                                        className="Pixso-paragraph-3_1573"
                                    >
                                        {"1000+"}
                                    </p>
                                    <p
                                        id="3_1574"
                                        className="Pixso-paragraph-3_1574"
                                    >
                                        {"тем в базе знаний"}
                                    </p>
                                </div>
                                <div id="3_1575" className="Pixso-frame-3_1575">
                                    <p
                                        id="3_1576"
                                        className="Pixso-paragraph-3_1576"
                                    >
                                        {"500+"}
                                    </p>
                                    <p
                                        id="3_1577"
                                        className="Pixso-paragraph-3_1577"
                                    >
                                        {"межпредметных связей"}
                                    </p>
                                </div>
                                <div id="3_1578" className="Pixso-frame-3_1578">
                                    <p
                                        id="3_1579"
                                        className="Pixso-paragraph-3_1579"
                                    >
                                        {"200+"}
                                    </p>
                                    <p
                                        id="3_1580"
                                        className="Pixso-paragraph-3_1580"
                                    >
                                        {"семей уже используют"}
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="3_1581" className="Pixso-frame-3_1581">
                    <div className="frame-content-3_1581">
                        <SectionHeader
                            id="7_237"
                            className="Pixso-instance-7_237"
                            slot_7_212={
                                <p id="11_86" className="Pixso-paragraph-11_86">
                                    {
                                        "Мы строим карту знаний, где каждая тема занимает своё место."
                                    }
                                </p>
                            }
                            slot_7_211={
                                <p id="11_85" className="Pixso-paragraph-11_85">
                                    {
                                        "Система, которая возвращает логику и смысл"
                                    }
                                </p>
                            }
                            slot_7_208={
                                <Eyebrow
                                    id="11_82"
                                    className="Pixso-instance-11_82"
                                    slot_4_273={
                                        <p
                                            id="11_84"
                                            className="Pixso-paragraph-11_84"
                                        >
                                            {"ОСНОВА"}
                                        </p>
                                    }
                                ></Eyebrow>
                            }
                        ></SectionHeader>
                        <div id="3_1585" className="Pixso-frame-3_1585">
                            <div className="frame-content-3_1585">
                                <div
                                    id="3_1586"
                                    className="stroke-wrapper-3_1586"
                                >
                                    <div className="Pixso-frame-3_1586">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1586">
                                            <div
                                                id="4_209"
                                                className="Pixso-vector-4_209"
                                            ></div>
                                            <p
                                                id="3_1587"
                                                className="Pixso-paragraph-3_1587"
                                            >
                                                {"ЦЕЛОСТНОСТЬ"}
                                            </p>
                                            <p
                                                id="3_1588"
                                                className="Pixso-paragraph-3_1588"
                                            >
                                                {"Целостность"}
                                            </p>
                                            <p
                                                id="3_1589"
                                                className="Pixso-paragraph-3_1589"
                                            >
                                                {
                                                    "Знания имеют смысл только в системе: одна тема вытекает из другой, предметы переплетаются. Мы строим карту, где эти связи видны, — чтобы ребёнок видел не отдельные куски, а живую картину мира."
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1586"></div>
                                </div>
                                <div
                                    id="3_1590"
                                    className="stroke-wrapper-3_1590"
                                >
                                    <div className="Pixso-frame-3_1590">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1590">
                                            <div
                                                id="4_215"
                                                className="Pixso-vector-4_215"
                                            ></div>
                                            <p
                                                id="3_1591"
                                                className="Pixso-paragraph-3_1591"
                                            >
                                                {"СВОБОДА"}
                                            </p>
                                            <p
                                                id="3_1592"
                                                className="Pixso-paragraph-3_1592"
                                            >
                                                {"Свобода без хаоса"}
                                            </p>
                                            <p
                                                id="3_1593"
                                                className="Pixso-paragraph-3_1593"
                                            >
                                                {
                                                    "Ребёнок идёт за своим интересом — в рамках маршрута, который ведёт к аттестации. Интерес и ответственность соединяются в едином плане. Мы не загоняем в рамки, но и не отпускаем без карты."
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1590"></div>
                                </div>
                                <div
                                    id="3_1594"
                                    className="stroke-wrapper-3_1594"
                                >
                                    <div className="Pixso-frame-3_1594">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1594">
                                            <div
                                                id="4_224"
                                                className="Pixso-vector-4_224"
                                            ></div>
                                            <p
                                                id="3_1595"
                                                className="Pixso-paragraph-3_1595"
                                            >
                                                {"СВЯЗИ"}
                                            </p>
                                            <p
                                                id="3_1596"
                                                className="Pixso-paragraph-3_1596"
                                            >
                                                {"Связь с жизнью"}
                                            </p>
                                            <p
                                                id="3_1597"
                                                className="Pixso-paragraph-3_1597"
                                            >
                                                {
                                                    "Каждая тема показывает, где она применяется: в профессии, в быту, в других предметах. Знания не висят в воздухе — они имеют смысл здесь и сейчас."
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1594"></div>
                                </div>
                                <div
                                    id="3_1598"
                                    className="stroke-wrapper-3_1598"
                                >
                                    <div className="Pixso-frame-3_1598">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1598">
                                            <div
                                                id="4_227"
                                                className="Pixso-vector-4_227"
                                            ></div>
                                            <p
                                                id="3_1599"
                                                className="Pixso-paragraph-3_1599"
                                            >
                                                {"СООБЩЕСТВО"}
                                            </p>
                                            <p
                                                id="3_1600"
                                                className="Pixso-paragraph-3_1600"
                                            >
                                                {"Сила в сообществе"}
                                            </p>
                                            <p
                                                id="3_1601"
                                                className="Pixso-paragraph-3_1601"
                                            >
                                                {
                                                    "Карта знаний — живая, её развивают методисты и родители вместе. Это общее дело: вы делитесь опытом, находите единомышленников и влияете на то, как учатся другие."
                                                }
                                            </p>
                                            <p
                                                id="8_13"
                                                className="Pixso-paragraph-8_13"
                                            >
                                                {"На платформе VEDO Hub →"}
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1598"></div>
                                </div>
                            </div>
                        </div>
                        <div id="3_1602" className="Pixso-frame-3_1602">
                            <div className="frame-content-3_1602">
                                <div
                                    id="3_1603"
                                    className="Pixso-vector-3_1603"
                                ></div>
                                <p
                                    id="3_1606"
                                    className="Pixso-paragraph-3_1606"
                                >
                                    {
                                        "«Дай пять» — это как экспедиция. У вас есть карта, где уже проложены маршруты. Но вы сами выбираете, по какой тропе идти. Рядом — те, кто прошёл этот путь раньше. И вы всегда знаете, где находитесь и что ждёт впереди."
                                    }
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="3_1607" className="Pixso-frame-3_1607">
                    <div className="frame-content-3_1607">
                        <SectionHeader
                            id="7_243"
                            className="Pixso-instance-7_243"
                            slot_7_212={
                                <p id="11_81" className="Pixso-paragraph-11_81">
                                    {
                                        "Отвечаем на сомнения до того, как они возникнут"
                                    }
                                </p>
                            }
                            slot_7_211={
                                <p id="11_80" className="Pixso-paragraph-11_80">
                                    {"Частые вопросы"}
                                </p>
                            }
                            slot_7_208={
                                <Eyebrow
                                    id="11_77"
                                    className="Pixso-instance-11_77"
                                    slot_4_273={
                                        <p
                                            id="11_79"
                                            className="Pixso-paragraph-11_79"
                                        >
                                            {"FAQ"}
                                        </p>
                                    }
                                ></Eyebrow>
                            }
                        ></SectionHeader>
                        <div id="3_1611" className="Pixso-frame-3_1611">
                            <div className="frame-content-3_1611">
                                <div
                                    id="3_1612"
                                    className="stroke-wrapper-3_1612"
                                >
                                    <div className="Pixso-frame-3_1612">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1612">
                                            <div
                                                id="4_291"
                                                className="Pixso-frame-4_291"
                                            >
                                                <div className="frame-content-4_291">
                                                    <p
                                                        id="3_1613"
                                                        className="Pixso-paragraph-3_1613"
                                                    >
                                                        {
                                                            "Что, если мы не успеем подготовиться к аттестации? Вдруг забудем какую-то тему?"
                                                        }
                                                    </p>
                                                    <div
                                                        id="4_274"
                                                        className="Pixso-vector-4_274"
                                                    ></div>
                                                </div>
                                            </div>
                                            <p
                                                id="3_1614"
                                                className="Pixso-paragraph-3_1614"
                                            >
                                                {
                                                    "Система заранее показывает пробелы: какие темы пройдены, какие остались, и сколько на них нужно времени. Вы видите покрытие ФГОС в реальном времени — ни одна тема не потеряется. Аттестация перестаёт быть стрессом, потому что вы точно знаете: всё под контролем."
                                                }
                                            </p>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1612"></div>
                                </div>
                                <div
                                    id="3_1615"
                                    className="stroke-wrapper-3_1615"
                                >
                                    <div className="Pixso-frame-3_1615">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1615">
                                            <div
                                                id="4_292"
                                                className="Pixso-frame-4_292"
                                            >
                                                <div className="frame-content-4_292">
                                                    <p
                                                        id="3_1616"
                                                        className="Pixso-paragraph-3_1616"
                                                    >
                                                        {
                                                            "Что если я сам не разбираюсь в предмете?"
                                                        }
                                                    </p>
                                                    <div
                                                        id="4_276"
                                                        className="Pixso-vector-4_276"
                                                    ></div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1615"></div>
                                </div>
                                <div
                                    id="3_1618"
                                    className="stroke-wrapper-3_1618"
                                >
                                    <div className="Pixso-frame-3_1618">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1618">
                                            <div
                                                id="4_293"
                                                className="Pixso-frame-4_293"
                                            >
                                                <div className="frame-content-4_293">
                                                    <p
                                                        id="3_1619"
                                                        className="Pixso-paragraph-3_1619"
                                                    >
                                                        {
                                                            "Можно ли вести несколько детей одновременно?"
                                                        }
                                                    </p>
                                                    <div
                                                        id="4_279"
                                                        className="Pixso-vector-4_279"
                                                    ></div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1618"></div>
                                </div>
                                <div
                                    id="3_1621"
                                    className="stroke-wrapper-3_1621"
                                >
                                    <div className="Pixso-frame-3_1621">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1621">
                                            <div
                                                id="4_294"
                                                className="Pixso-frame-4_294"
                                            >
                                                <div className="frame-content-4_294">
                                                    <p
                                                        id="3_1622"
                                                        className="Pixso-paragraph-3_1622"
                                                    >
                                                        {"Это бесплатно?"}
                                                    </p>
                                                    <div
                                                        id="4_282"
                                                        className="Pixso-vector-4_282"
                                                    ></div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1621"></div>
                                </div>
                                <div
                                    id="3_1624"
                                    className="stroke-wrapper-3_1624"
                                >
                                    <div className="Pixso-frame-3_1624">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1624">
                                            <div
                                                id="4_295"
                                                className="Pixso-frame-4_295"
                                            >
                                                <div className="frame-content-4_295">
                                                    <p
                                                        id="3_1625"
                                                        className="Pixso-paragraph-3_1625"
                                                    >
                                                        {
                                                            "Сколько времени займёт настройка?"
                                                        }
                                                    </p>
                                                    <div
                                                        id="4_285"
                                                        className="Pixso-vector-4_285"
                                                    ></div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1624"></div>
                                </div>
                                <div
                                    id="3_1627"
                                    className="stroke-wrapper-3_1627"
                                >
                                    <div className="Pixso-frame-3_1627">
                                        <div className="shadow-blend-unknown-1"></div>
                                        <div className="shadow-blend-unknown-0"></div>
                                        <div className="frame-content-3_1627">
                                            <div
                                                id="4_296"
                                                className="Pixso-frame-4_296"
                                            >
                                                <div className="frame-content-4_296">
                                                    <p
                                                        id="3_1628"
                                                        className="Pixso-paragraph-3_1628"
                                                    >
                                                        {
                                                            "Мои данные в безопасности?"
                                                        }
                                                    </p>
                                                    <div
                                                        id="4_288"
                                                        className="Pixso-vector-4_288"
                                                    ></div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="stroke-3_1627"></div>
                                </div>
                                <div id="6_111" className="Pixso-frame-6_111">
                                    <div className="frame-content-6_111">
                                        <div
                                            id="6_112"
                                            className="Pixso-frame-6_112"
                                        >
                                            <div className="frame-content-6_112">
                                                <p
                                                    id="6_113"
                                                    className="Pixso-paragraph-6_113"
                                                >
                                                    {"Кто создаёт граф знаний?"}
                                                </p>
                                                <div
                                                    id="6_114"
                                                    className="Pixso-vector-6_114"
                                                ></div>
                                            </div>
                                        </div>
                                        <p
                                            id="6_115"
                                            className="Pixso-paragraph-6_115"
                                        >
                                            {
                                                "Граф знаний создаётся методистами совместно с родительским сообществом на основе онтологии VEDO Hub. Родители могут предлагать связи между темами, делиться опытом и учебными находками. Это живая система: граф постоянно пополняется и уточняется — академическая точность методистов встречается с реальным родительским опытом."
                                            }
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="3_1630" className="Pixso-frame-3_1630">
                    <div className="frame-content-3_1630">
                        <p id="3_1631" className="Pixso-paragraph-3_1631">
                            {"Готовы построить маршрут?"}
                        </p>
                        <p id="3_1632" className="Pixso-paragraph-3_1632">
                            {
                                "Бесплатно для всей семьи. 5 минут на маршрут. Доплата — только за ИИ-ассистента. 200+ семей."
                            }
                        </p>
                        <div id="4_428" className="Pixso-frame-4_428">
                            <div
                                id="4_429"
                                className="Pixso-vector-4_429"
                            ></div>
                            <p id="4_433" className="Pixso-paragraph-4_433">
                                {"Не проходите курс — достигайте цель"}
                            </p>
                        </div>
                        <div id="3_1633" className="Pixso-frame-3_1633">
                            <div
                                id="3_1634"
                                className="Pixso-vector-3_1634"
                            ></div>
                            <p id="3_1639" className="Pixso-paragraph-3_1639">
                                {"Начать строить маршрут"}
                            </p>
                            <div id="4_91" className="Pixso-frame-4_91">
                                <div className="frame-content-4_91">
                                    <div
                                        id="4_92"
                                        className="Pixso-vector-4_92"
                                    ></div>
                                </div>
                            </div>
                        </div>
                        <div id="4_89" className="Pixso-vector-4_89"></div>
                    </div>
                </div>
                <div id="3_1640" className="Pixso-frame-3_1640">
                    <div className="frame-content-3_1640">
                        <div id="3_1641" className="Pixso-frame-3_1641">
                            <div className="frame-content-3_1641">
                                <div id="3_1642" className="Pixso-frame-3_1642">
                                    <div
                                        id="3_1643"
                                        className="Pixso-vector-3_1643"
                                    ></div>
                                    <p
                                        id="3_1648"
                                        className="Pixso-paragraph-3_1648"
                                    >
                                        {"Дай пять"}
                                    </p>
                                </div>
                                <div id="3_1649" className="Pixso-frame-3_1649">
                                    <div
                                        id="3_1650"
                                        className="Pixso-frame-3_1650"
                                    >
                                        <p
                                            id="3_1651"
                                            className="Pixso-paragraph-3_1651"
                                        >
                                            {"Продукт"}
                                        </p>
                                        <p
                                            id="3_1652"
                                            className="Pixso-paragraph-3_1652"
                                        >
                                            {"Как это работает"}
                                        </p>
                                        <p
                                            id="6_121"
                                            className="Pixso-paragraph-6_121"
                                        >
                                            {"Тарифы"}
                                        </p>
                                        <p
                                            id="3_1653"
                                            className="Pixso-paragraph-3_1653"
                                        >
                                            {"Возможности"}
                                        </p>
                                        <p
                                            id="3_1654"
                                            className="Pixso-paragraph-3_1654"
                                        >
                                            {"Отзывы"}
                                        </p>
                                    </div>
                                    <div
                                        id="3_1655"
                                        className="Pixso-frame-3_1655"
                                    >
                                        <p
                                            id="3_1656"
                                            className="Pixso-paragraph-3_1656"
                                        >
                                            {"Поддержка"}
                                        </p>
                                        <p
                                            id="3_1657"
                                            className="Pixso-paragraph-3_1657"
                                        >
                                            {"FAQ"}
                                        </p>
                                        <p
                                            id="3_1658"
                                            className="Pixso-paragraph-3_1658"
                                        >
                                            {"Контакты"}
                                        </p>
                                        <p
                                            id="3_1659"
                                            className="Pixso-paragraph-3_1659"
                                        >
                                            {"Блог"}
                                        </p>
                                    </div>
                                    <div
                                        id="3_1660"
                                        className="Pixso-frame-3_1660"
                                    >
                                        <p
                                            id="3_1661"
                                            className="Pixso-paragraph-3_1661"
                                        >
                                            {"Правовая информация"}
                                        </p>
                                        <p
                                            id="3_1662"
                                            className="Pixso-paragraph-3_1662"
                                        >
                                            {"Политика конфиденциальности"}
                                        </p>
                                        <p
                                            id="3_1663"
                                            className="Pixso-paragraph-3_1663"
                                        >
                                            {"Условия использования"}
                                        </p>
                                        <p
                                            id="3_1664"
                                            className="Pixso-paragraph-3_1664"
                                        >
                                            {"Согласие на обработку данных"}
                                        </p>
                                    </div>
                                    <div
                                        id="7_318"
                                        className="Pixso-frame-7_318"
                                    >
                                        <p
                                            id="7_319"
                                            className="Pixso-paragraph-7_319"
                                        >
                                            {"Для бизнеса"}
                                        </p>
                                        <p
                                            id="7_320"
                                            className="Pixso-paragraph-7_320"
                                        >
                                            {"Для бизнеса"}
                                        </p>
                                        <p
                                            id="8_12"
                                            className="Pixso-paragraph-8_12"
                                        >
                                            {"VEDO Hub"}
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div id="3_1665" className="Pixso-frame-3_1665">
                            <div className="frame-content-3_1665">
                                <p
                                    id="3_1666"
                                    className="Pixso-paragraph-3_1666"
                                >
                                    {"© 2026 Дай пять. Все права защищены."}
                                </p>
                                <div id="3_1667" className="Pixso-frame-3_1667">
                                    <div
                                        id="3_1668"
                                        className="Pixso-vector-3_1668"
                                    ></div>
                                    <div
                                        id="3_1671"
                                        className="Pixso-vector-3_1671"
                                    ></div>
                                    <div
                                        id="3_1676"
                                        className="Pixso-vector-3_1676"
                                    ></div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}