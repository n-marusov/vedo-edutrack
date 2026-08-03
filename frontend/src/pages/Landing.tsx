import { Link } from 'react-router';
import { appConfig } from '../config';
import { Card } from '../shared/components/Card';

const values = [
  {
    title: 'Knowledge Graph Routes',
    text: 'Personalized learning paths built on the knowledge graph: each step follows strict prerequisites toward your goal.',
  },
  {
    title: 'Progress Tracking',
    text: 'Plan vs actual comparison with deviation alerts and readiness forecasts — always know where you stand.',
  },
  {
    title: 'Gap Diagnosis',
    text: 'Find the root cause of learning gaps by climbing prerequisite links to the first unmastered module.',
  },
];

/** Landing page — public marketing page at "/". */
export function LandingPage() {
  return (
    <div className="min-h-screen">
      <header className="flex items-center justify-between px-8 py-4">
        <span className="text-xl font-bold text-blue-600">VEDO EduTrack</span>
        <Link
          to="/login"
          className="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          Sign In
        </Link>
      </header>

      <section className="bg-gradient-to-b from-blue-50 to-white px-8 py-24 text-center">
        <h1 className="mx-auto max-w-3xl text-5xl font-bold tracking-tight text-gray-900">
          VEDO EduTrack
        </h1>
        <p className="mx-auto mt-6 max-w-2xl text-xl text-gray-600">
          Персональная образовательная траектория на основе графа знаний
        </p>
        <div className="mt-10 flex items-center justify-center gap-4">
          <Link
            to="/login"
            className="rounded-md bg-blue-600 px-6 py-3 text-base font-medium text-white hover:bg-blue-700"
          >
            Get started
          </Link>
          <a
            href="#features"
            className="rounded-md bg-white px-6 py-3 text-base font-medium text-gray-700 ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
          >
            Learn more
          </a>
        </div>
      </section>

      <section id="features" className="mx-auto max-w-6xl px-8 py-16">
        <h2 className="text-center text-2xl font-semibold text-gray-900">Why VEDO EduTrack</h2>
        <div className="mt-10 grid gap-6 md:grid-cols-3">
          {values.map((v) => (
            <Card key={v.title} title={v.title}>
              <p className="text-sm text-gray-600">{v.text}</p>
            </Card>
          ))}
        </div>
      </section>

      <footer className="border-t border-gray-100 px-8 py-8 text-center text-sm text-gray-500">
        © {new Date().getFullYear()} VEDO EduTrack · {appConfig.env} · v{appConfig.version}
      </footer>
    </div>
  );
}
