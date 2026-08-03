import type { FC } from 'react';

/**
 * VEDO EduTrack root application component.
 *
 * TODO: Add routing, layout, and feature modules (M0.3+).
 */
export const App: FC = () => {
  return (
    <main className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900">VEDO EduTrack</h1>
        <p className="mt-4 text-lg text-gray-600">Сервис образовательных маршрутов</p>
      </div>
    </main>
  );
};
