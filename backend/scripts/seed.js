// MongoDB seed script - runs on container initialization
// Creates test user, project, and tasks

db = db.getSiblingDB('project_management');

// Create test user with password: testpass123
// Password hash for "testpass123" using bcrypt cost 12
const testUser = {
  _id: 'test-user-001',
  email: 'test@example.com',
  password_hash: '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIpSVeu1Eq',
  display_name: 'Test User',
  is_active: true,
  created_at: new Date(),
  updated_at: new Date()
};

db.users.insertOne(testUser);
print('✓ Created test user: test@example.com / testpass123');

// Create default workflow
const workflow = {
  _id: 'workflow-001',
  project_id: 'project-001',
  name: 'Default Workflow',
  statuses: [
    {
      id: 'todo',
      name: 'To Do',
      category: 'todo',
      order: 0
    },
    {
      id: 'in-progress',
      name: 'In Progress',
      category: 'in_progress',
      order: 1
    },
    {
      id: 'done',
      name: 'Done',
      category: 'done',
      order: 2
    }
  ],
  transitions: [
    { from: 'todo', to: 'in-progress' },
    { from: 'in-progress', to: 'done' },
    { from: 'in-progress', to: 'todo' },
    { from: 'done', to: 'in-progress' }
  ],
  created_at: new Date(),
  updated_at: new Date()
};

db.workflows.insertOne(workflow);
print('✓ Created default workflow');

// Create test project
const project = {
  _id: 'project-001',
  name: 'Demo Project',
  key: 'DEMO',
  description: 'A demo project with sample tasks',
  workflow_id: 'workflow-001',
  custom_fields: [],
  members: ['test-user-001'],
  created_at: new Date(),
  updated_at: new Date()
};

db.projects.insertOne(project);
print('✓ Created demo project');

// Create default sprint
const sprint = {
  _id: 'sprint-001',
  project_id: 'project-001',
  name: 'Sprint 1',
  goal: 'Initial sprint with demo tasks',
  status: 'active',
  is_default: true,
  completed_points: 0,
  total_points: 13,
  velocity: 0,
  issue_count: 3,
  complete_issue_count: 0,
  created_at: new Date(),
  updated_at: new Date()
};

db.sprints.insertOne(sprint);
print('✓ Created default sprint');

// Create 3 sample tasks with different statuses
const tasks = [
  {
    _id: 'issue-001',
    issue_key: 'DEMO-1',
    project_id: 'project-001',
    type: 'task',
    title: 'Setup project infrastructure',
    description: 'Configure Docker, database, and initial project structure',
    status: 'done',
    priority: 'high',
    assignee_id: 'test-user-001',
    reporter_id: 'test-user-001',
    sprint_id: 'sprint-001',
    labels: ['infrastructure', 'setup'],
    story_points: 5,
    custom_fields: {},
    watchers: ['test-user-001'],
    version: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    _id: 'issue-002',
    issue_key: 'DEMO-2',
    project_id: 'project-001',
    type: 'task',
    title: 'Implement authentication system',
    description: 'Add JWT-based authentication with login and registration',
    status: 'in-progress',
    priority: 'highest',
    assignee_id: 'test-user-001',
    reporter_id: 'test-user-001',
    sprint_id: 'sprint-001',
    labels: ['security', 'auth'],
    story_points: 8,
    custom_fields: {},
    watchers: ['test-user-001'],
    version: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    _id: 'issue-003',
    issue_key: 'DEMO-3',
    project_id: 'project-001',
    type: 'bug',
    title: 'Fix responsive layout on mobile',
    description: 'Ensure all pages work correctly at 375px width',
    status: 'todo',
    priority: 'medium',
    assignee_id: 'test-user-001',
    reporter_id: 'test-user-001',
    sprint_id: 'sprint-001',
    labels: ['ui', 'mobile'],
    story_points: 3,
    custom_fields: {},
    watchers: ['test-user-001'],
    version: 1,
    created_at: new Date(),
    updated_at: new Date()
  }
];

db.issues.insertMany(tasks);
print('✓ Created 3 sample tasks');

// Create indexes
db.users.createIndex({ email: 1 }, { unique: true });
db.projects.createIndex({ key: 1 }, { unique: true });
db.issues.createIndex({ issue_key: 1 }, { unique: true });
db.issues.createIndex({ project_id: 1 });
db.issues.createIndex({ assignee_id: 1 });
db.sprints.createIndex({ project_id: 1 });
db.workflows.createIndex({ project_id: 1 });

print('✓ Created database indexes');
print('\n=== Seed completed successfully ===');
print('Test credentials:');
print('  Email: test@example.com');
print('  Password: testpass123');
