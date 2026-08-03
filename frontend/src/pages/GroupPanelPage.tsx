import { useNavigate } from 'react-router';
import { GroupPanel, type LearnerCard } from '../features/visualization';
import { RoleGate } from '../shared/guards/RoleGate';

const demoLearners: LearnerCard[] = [
  {
    id: 'demo-misha',
    name: 'Миша',
    currentModule: 'percent',
    fgosCoverage: 70,
    forecastStatus: 'at-risk',
    attentionFlag: true,
  },
  {
    id: 'demo-katya',
    name: 'Катя',
    currentModule: 'fractions',
    fgosCoverage: 90,
    forecastStatus: 'on-track',
    attentionFlag: false,
  },
  {
    id: 'demo-anya',
    name: 'Аня',
    currentModule: 'mechanics',
    fgosCoverage: 60,
    forecastStatus: 'not-on-track',
    attentionFlag: true,
  },
];

/** GroupPanelPage — role-aware learner group panel (M2 F4.7). */
export function GroupPanelPage() {
  const navigate = useNavigate();
  return (
    <RoleGate requiredRole={['parent', 'teacher', 'methodologist', 'admin']}>
      <GroupPanel
        title="Мои ученики"
        learners={demoLearners}
        onSelect={(learnerId) => navigate(`/dashboard/learner?learner=${learnerId}`)}
      />
    </RoleGate>
  );
}
