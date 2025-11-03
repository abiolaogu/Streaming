import { RoadmapIssue } from '../types';

export const roadmapIssues: RoadmapIssue[] = [
  // Phase 1: Foundation & Infrastructure
  {
    id: 'ISSUE-001',
    title: 'Project Initialization & Repository Setup',
    phase: 'Phase 1: Foundation & Infrastructure',
    // FIX: Add 'as const' to prevent type widening to string.
    status: 'Done' as const,
    priority: 'P0',
    effort: '2h',
    owner: 'Platform Architect',
    description: 'Initialize monorepo structure, configure build tools, set up workspace management.',
  },
  {
    id: 'ISSUE-002',
    title: 'Docker Infrastructure & Local Development Setup',
    phase: 'Phase 1: Foundation & Infrastructure',
    // FIX: Add 'as const' to prevent type widening to string.
    status: 'In Progress' as const,
    priority: 'P0',
    effort: '4h',
    owner: 'DevOps Lead',
    description: 'Create docker-compose.yml for local development, base Dockerfiles for Go/Python/Node.',
  },
  {
    id: 'ISSUE-003',
    title: 'YugabyteDB Schema & Initialization',
    phase: 'Phase 1: Foundation & Infrastructure',
    // FIX: Add 'as const' to prevent type widening to string.
    status: 'In Progress' as const,
    priority: 'P0',
    effort: '6h',
    owner: 'Database Lead',
    description: 'Design and implement complete YugabyteDB (YSQL) schema for all services. Replication factor 3 across all PoPs.',
  },
  {
    id: 'ISSUE-004',
    title: 'Message Broker Setup (Kafka/RabbitMQ)',
    phase: 'Phase 1: Foundation & Infrastructure',
    // FIX: Add 'as const' to prevent type widening to string.
    status: 'Not Started' as const,
    priority: 'P0',
    effort: '4h',
    owner: 'Platform Architect',
    description: 'Configure Kafka/RabbitMQ, define topics, implement schema registry.',
  },
  {
    id: 'ISSUE-005',
    title: 'API Gateway Configuration (Kong)',
    phase: 'Phase 1: Foundation & Infrastructure',
    // FIX: Add 'as const' to prevent type widening to string.
    status: 'Not Started' as const,
    priority: 'P0',
    effort: '3h',
    owner: 'Platform Architect',
    description: 'Set up Kong API Gateway with rate limiting, authentication, routing.',
  },
  {
    id: 'ISSUE-006',
    title: 'Shared Go Libraries (common-go)',
    phase: 'Phase 1: Foundation & Infrastructure',
    // FIX: Add 'as const' to prevent type widening to string.
    status: 'In Progress' as const,
    priority: 'P0',
    effort: '6h',
    owner: 'Backend Lead',
    description: 'Create reusable Go packages for logging, middleware, authentication, error handling.',
  },
  {
    id: 'ISSUE-007',
    title: 'Shared TypeScript/Node Libraries (common-ts)',
    phase: 'Phase 1: Foundation & Infrastructure',
    // FIX: Add 'as const' to prevent type widening to string.
    status: 'In Progress' as const,
    priority: 'P0',
    effort: '4h',
    owner: 'Frontend Lead',
    description: 'Create reusable Node.js/TypeScript packages for shared utilities.',
  },
  {
    id: 'ISSUE-008',
    title: 'Shared Python Libraries',
    phase: 'Phase 1: Foundation & Infrastructure',
    // FIX: Add 'as const' to prevent type widening to string.
    status: 'Not Started' as const,
    priority: 'P0',
    effort: '3h',
    owner: 'Data Lead',
    description: 'Create reusable Python packages for ML, analytics, logging.',
  },
  {
    id: 'ISSUE-009',
    title: 'Protocol Buffers Definitions (proto)',
    phase: 'Phase 1: Foundation & Infrastructure',
    // FIX: Add 'as const' to prevent type widening to string.
    status: 'Not Started' as const,
    priority: 'P0',
    effort: '3h',
    owner: 'Platform Architect',
    description: 'Define Protocol Buffer schemas for service communication.',
  },
  {
    id: 'ISSUE-010',
    title: 'CI/CD Pipeline Foundation (GitHub Actions)',
    phase: 'Phase 1: Foundation & Infrastructure',
    // FIX: Add 'as const' to prevent type widening to string.
    status: 'Not Started' as const,
    priority: 'P0',
    effort: '5h',
    owner: 'DevOps Lead',
    description: 'Set up GitHub Actions workflows for lint, test, build, security scan.',
  },
  // Phase 2: Core Microservices - Backend
  ...Array.from({ length: 25 }, (_, i) => ({
    id: `ISSUE-${String(i + 11).padStart(3, '0')}`,
    title: 'TBD',
    phase: 'Phase 2: Core Microservices - Backend',
    status: 'Not Started' as const,
    priority: 'P1',
    effort: 'TBD',
    owner: 'TBD',
    description: 'This issue is pending detailed planning.',
  })),
  // Phase 3: Frontend Applications
    ...Array.from({ length: 5 }, (_, i) => ({
    id: `ISSUE-${String(i + 31).padStart(3, '0')}`,
    title: 'TBD',
    phase: 'Phase 3: Frontend Applications',
    status: 'Not Started' as const,
    priority: 'P1',
    effort: 'TBD',
    owner: 'TBD',
    description: 'This issue is pending detailed planning.',
  })),
  // Phase 4: Infrastructure & Deployment
  ...Array.from({ length: 10 }, (_, i) => ({
    id: `ISSUE-${String(i + 36).padStart(3, '0')}`,
    title: 'TBD',
    phase: 'Phase 4: Infrastructure & Deployment',
    status: 'Not Started' as const,
    priority: 'P1',
    effort: 'TBD',
    owner: 'TBD',
    description: 'This issue is pending detailed planning.',
  })),
  // Phase 5: Testing & Quality Assurance
  ...Array.from({ length: 4 }, (_, i) => ({
    id: `ISSUE-${String(i + 46).padStart(3, '0')}`,
    title: 'TBD',
    phase: 'Phase 5: Testing & Quality Assurance',
    status: 'Not Started' as const,
    priority: 'P1',
    effort: 'TBD',
    owner: 'TBD',
    description: 'This issue is pending detailed planning.',
  })),
  // Phase 6: Launch & Optimization
  {
    id: 'ISSUE-050',
    title: 'Staging Environment Validation',
    phase: 'Phase 6: Launch & Optimization',
    // FIX: Add 'as const' to prevent type widening to string.
    status: 'Not Started' as const,
    priority: 'P1',
    effort: '8h',
    owner: 'QA Lead',
    description: 'Validate complete platform in staging environment.',
  },
  {
    id: 'ISSUE-051',
    title: 'Production Launch & Go-Live Support',
    phase: 'Phase 6: Launch & Optimization',
    // FIX: Add 'as const' to prevent type widening to string.
    status: 'Not Started' as const,
    priority: 'P0',
    effort: 'Ongoing',
    owner: 'Platform Lead',
    description: 'Launch platform to production and provide support.',
  },
].map(issue => {
    // This is where we can manually update the details from the markdown
    const details: Record<string, Partial<RoadmapIssue>> = {
        'ISSUE-011': { title: 'Auth Service - Core Implementation', owner: 'Backend Lead', effort: '8h', description: 'Implement authentication service with JWT, OAuth2.0, MFA support.' },
        'ISSUE-012': { title: 'User Service - Profile & Preferences', owner: 'Backend Lead', effort: '6h', status: 'In Progress' as const, description: 'Implement user profile management, preferences, device management.' },
        'ISSUE-013': { title: 'Content Service - Metadata & Catalog', owner: 'Backend Lead', effort: '8h', description: 'Implement content catalog with metadata, search, and recommendations hooks.' },
        'ISSUE-014': { title: 'Streaming Service - Manifest Generation & Token', owner: 'Backend Lead', effort: '10h', description: 'Implement streaming manifest generation with tokenization and ABR selection.' },
        'ISSUE-028': { title: 'Video Player SDK Development', owner: 'Frontend Lead', effort: '12h', status: 'In Progress' as const, description: 'Develop cross-platform video player SDK with ABR, DRM, analytics.' },
        'ISSUE-031': { title: 'Web App - Next.js Setup & Authentication', owner: 'Frontend Lead', effort: '8h', status: 'In Progress' as const, description: 'Set up Next.js web application with authentication and core pages.' },
        'ISSUE-034': { title: 'Admin Dashboard - Next.js Web App', owner: 'Frontend Lead', effort: '10h', status: 'In Progress' as const, description: 'Build admin dashboard for content management and analytics.' },
    };
    return { ...issue, ...details[issue.id] };
})
];