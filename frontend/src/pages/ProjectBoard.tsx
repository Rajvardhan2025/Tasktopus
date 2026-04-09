import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
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
import { Wifi, WifiOff, LayoutGrid, Zap, GitBranch } from 'lucide-react';

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
    return <div className="flex items-center justify-center h-64">Loading project...</div>;
  }

  if (!project) {
    return <div className="flex items-center justify-center h-64">Project not found</div>;
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex-1">
          <div className="flex items-center gap-3">
            <h1 className="text-3xl font-bold">{project.name}</h1>
            <Badge variant="outline" className="font-mono">{project.key}</Badge>
            <Badge variant={isConnected ? 'default' : 'secondary'} className="gap-1">
              {isConnected ? (
                <>
                  <Wifi className="w-3 h-3" />
                  Live
                </>
              ) : (
                <>
                  <WifiOff className="w-3 h-3" />
                  Offline
                </>
              )}
            </Badge>
          </div>
          {project.description && (
            <p className="text-muted-foreground mt-1">{project.description}</p>
          )}
        </div>
        <div className="flex items-center gap-3">
          {projectId && <ProjectMembersManager projectId={projectId} />}
          <SearchBar projectId={projectId} onSelectIssue={setSelectedIssue} />
        </div>
      </div>

      {/* View Tabs */}
      <div className="flex gap-2 mb-6 border-b">
        <Button
          variant={activeView === 'board' ? 'default' : 'ghost'}
          onClick={() => setActiveView('board')}
          className="rounded-b-none"
        >
          <LayoutGrid className="w-4 h-4 mr-2" />
          Board
        </Button>
        <Button
          variant={activeView === 'sprints' ? 'default' : 'ghost'}
          onClick={() => setActiveView('sprints')}
          className="rounded-b-none"
        >
          <Zap className="w-4 h-4 mr-2" />
          Sprints
        </Button>
        <Button
          variant={activeView === 'workflow' ? 'default' : 'ghost'}
          onClick={() => setActiveView('workflow')}
          className="rounded-b-none"
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
