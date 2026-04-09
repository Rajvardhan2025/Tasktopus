import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { projectsApi } from '@/lib/api';
import type { Issue } from '@/lib/api';
import { KanbanBoard } from '@/components/KanbanBoard';
import { SprintManagement } from '@/components/SprintManagement';
import { WorkflowManagement } from '@/components/WorkflowManagement';
import { SearchBar } from '@/components/SearchBar';
import { ProjectMembersManager } from '@/components/ProjectMembersManager';
import { IssueDialog } from '@/components/IssueDialog';
import { useWebSocket } from '@/lib/websocket';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Wifi, WifiOff, LayoutGrid, Zap, GitBranch, ChevronRight, FolderKanban } from 'lucide-react';

export function ProjectBoard() {
  const { projectId } = useParams<{ projectId: string }>();
  const [activeView, setActiveView] = useState<'board' | 'sprints' | 'workflow'>('board');
  const [selectedIssue, setSelectedIssue] = useState<Issue | null>(null);

  const { data, isLoading } = useQuery({
    queryKey: ['project', projectId],
    queryFn: () => projectsApi.get(projectId!),
    enabled: !!projectId,
  });

  const { isConnected, lastEvent } = useWebSocket(projectId || null);

  const queryClient = useQueryClient();

  useEffect(() => {
    if (!lastEvent || !projectId) return;

    if (
      lastEvent.type === 'issue_created' ||
      lastEvent.type === 'issue_updated' ||
      lastEvent.type === 'issue_moved' ||
      lastEvent.type === 'comment_added'
    ) {
      queryClient.invalidateQueries({ queryKey: ['issues', projectId] });
    }

    if (lastEvent.type === 'sprint_updated') {
      queryClient.invalidateQueries({ queryKey: ['sprints', projectId] });
    }

    queryClient.invalidateQueries({ queryKey: ['activity'] });
  }, [lastEvent, projectId, queryClient]);

  const project = data?.data?.data;

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center space-y-4">
          <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto"></div>
          <p className="text-muted-foreground">Loading project...</p>
        </div>
      </div>
    );
  }

  if (!project) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center space-y-4">
          <FolderKanban className="w-16 h-16 mx-auto text-muted-foreground" />
          <div>
            <h3 className="text-lg font-semibold mb-2">Project not found</h3>
            <p className="text-muted-foreground mb-4">The project you're looking for doesn't exist</p>
            <Button asChild>
              <Link to="/projects">Back to Projects</Link>
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="animate-fade-in">
      {/* Breadcrumbs */}
      <nav className="flex items-center gap-2 text-sm text-muted-foreground mb-4">
        <Link to="/projects" className="hover:text-foreground transition-colors">
          Projects
        </Link>
        <ChevronRight className="w-4 h-4" />
        <span className="text-foreground font-medium">{project.name}</span>
      </nav>

      {/* Project Header */}
      <div className="flex flex-col lg:flex-row items-start lg:items-center justify-between mb-6 gap-4">
        <div className="flex-1">
          <div className="flex flex-wrap items-center gap-3">
            <h1 className="text-3xl font-bold">{project.name}</h1>
            <Badge variant="outline" className="font-mono">{project.key}</Badge>
          </div>
          {project.description && (
            <p className="text-muted-foreground mt-2">{project.description}</p>
          )}
        </div>
        <div className="flex items-center gap-3 w-full lg:w-auto">
          {projectId && <ProjectMembersManager projectId={projectId} />}
          <SearchBar projectId={projectId} onSelectIssue={setSelectedIssue} />
        </div>
      </div>

      {/* View Tabs */}
      <div className="flex gap-2 mb-6 border-b overflow-x-auto">
        <Button
          variant={activeView === 'board' ? 'default' : 'ghost'}
          onClick={() => setActiveView('board')}
          className="rounded-b-none whitespace-nowrap"
        >
          <LayoutGrid className="w-4 h-4 mr-2" />
          Board
        </Button>
        <Button
          variant={activeView === 'sprints' ? 'default' : 'ghost'}
          onClick={() => setActiveView('sprints')}
          className="rounded-b-none whitespace-nowrap"
        >
          <Zap className="w-4 h-4 mr-2" />
          Sprints
        </Button>
        <Button
          variant={activeView === 'workflow' ? 'default' : 'ghost'}
          onClick={() => setActiveView('workflow')}
          className="rounded-b-none whitespace-nowrap"
        >
          <GitBranch className="w-4 h-4 mr-2" />
          Workflow
        </Button>
      </div>

      {/* Content */}
      {activeView === 'board' && <KanbanBoard projectId={projectId!} />}
      {activeView === 'sprints' && <SprintManagement projectId={projectId!} />}
      {activeView === 'workflow' && <WorkflowManagement projectId={projectId!} />}

      {/* Search Result Dialog */}
      {selectedIssue && (
        <IssueDialog
          issue={selectedIssue}
          open={!!selectedIssue}
          onClose={() => setSelectedIssue(null)}
        />
      )}
    </div>
  );
}
