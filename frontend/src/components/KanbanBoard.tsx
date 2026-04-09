import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { issuesApi } from '@/lib/api';
import type { Issue } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Plus } from 'lucide-react';
import { IssueDialog } from './IssueDialog';
import { CreateIssueDialog } from './CreateIssueDialog';
import api from '@/lib/api';
import { useToast } from '@/hooks/use-toast';

interface KanbanBoardProps {
  projectId: string;
}

interface Workflow {
  id: string;
  name: string;
  project_id: string;
  statuses: string[];
  transitions: Array<{
    from: string;
    to: string;
    conditions: Array<{ field: string; operator: string; value: string }>;
    actions: Array<{ type: string; params: Record<string, any> }>;
  }>;
}

const STATUS_LABELS: Record<string, string> = {
  to_do: 'To Do',
  in_progress: 'In Progress',
  in_review: 'In Review',
  done: 'Done',
};

const PRIORITY_COLORS: Record<string, string> = {
  lowest: 'bg-gray-200 text-gray-800',
  low: 'bg-blue-200 text-blue-800',
  medium: 'bg-yellow-200 text-yellow-800',
  high: 'bg-orange-200 text-orange-800',
  highest: 'bg-red-200 text-red-800',
};

export function KanbanBoard({ projectId }: KanbanBoardProps) {
  const [selectedIssue, setSelectedIssue] = useState<Issue | null>(null);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const queryClient = useQueryClient();
  const { toast } = useToast();

  const { data, isLoading } = useQuery({
    queryKey: ['issues', projectId],
    queryFn: () => issuesApi.list(projectId),
  });

  const { data: workflowData } = useQuery({
    queryKey: ['workflow', projectId],
    queryFn: async () => {
      const response = await api.get<{ data: Workflow }>(`/projects/${projectId}/workflow`);
      return response.data;
    },
  });

  const transitionMutation = useMutation({
    mutationFn: ({ issueId, toStatus, version }: { issueId: string; toStatus: string; version: number }) =>
      issuesApi.transition(issueId, toStatus, version),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['issues', projectId] });
      toast({
        title: 'Success',
        description: 'Issue status updated successfully',
      });
    },
    onError: (error: any) => {
      const message =
        error.response?.data?.error?.message ||
        error.response?.data?.message ||
        'Failed to update issue status';
      toast({
        title: 'Error',
        description: message,
        variant: 'destructive',
      });
    },
  });

  const issues = data?.data?.data?.issues || [];
  const workflow = workflowData?.data;
  const statuses = workflow?.statuses || ['to_do', 'in_progress', 'in_review', 'done'];

  const issuesByStatus = statuses.reduce((acc, status) => {
    acc[status] = issues.filter((issue) => issue.status === status);
    return acc;
  }, {} as Record<string, Issue[]>);

  const isTransitionAllowed = (fromStatus: string, toStatus: string): boolean => {
    if (!workflow) return true;
    return workflow.transitions.some((t) => t.from === fromStatus && t.to === toStatus);
  };

  const handleDragStart = (e: React.DragEvent, issue: Issue) => {
    e.dataTransfer.setData('issueId', issue.id);
    e.dataTransfer.setData('version', issue.version.toString());
    e.dataTransfer.setData('fromStatus', issue.status);
  };

  const handleDrop = (e: React.DragEvent, toStatus: string) => {
    e.preventDefault();
    const issueId = e.dataTransfer.getData('issueId');
    const version = parseInt(e.dataTransfer.getData('version'));
    const fromStatus = e.dataTransfer.getData('fromStatus');

    if (fromStatus === toStatus) return;

    if (!isTransitionAllowed(fromStatus, toStatus)) {
      toast({
        title: 'Invalid Transition',
        description: `Cannot move from ${STATUS_LABELS[fromStatus]} to ${STATUS_LABELS[toStatus]}`,
        variant: 'destructive',
      });
      return;
    }

    transitionMutation.mutate({ issueId, toStatus, version });
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
  };

  if (isLoading) {
    return <div className="flex items-center justify-center h-64">Loading...</div>;
  }

  return (
    <div className="h-full">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">Board</h2>
        <Button onClick={() => setCreateDialogOpen(true)}>
          <Plus className="w-4 h-4 mr-2" />
          Create Issue
        </Button>
      </div>

      <div className="overflow-x-auto pb-2">
        <div className="grid min-w-[980px] grid-cols-4 gap-4 h-[calc(100vh-220px)]">
          {statuses.map((status) => (
            <div
              key={status}
              className="flex flex-col bg-muted/50 rounded-lg p-4"
              onDrop={(e) => handleDrop(e, status)}
              onDragOver={handleDragOver}
            >
              <div className="flex items-center justify-between mb-4">
                <h3 className="font-semibold text-sm uppercase tracking-wide">
                  {STATUS_LABELS[status] || status}
                </h3>
                <Badge variant="secondary">{issuesByStatus[status]?.length || 0}</Badge>
              </div>

              <div className="flex-1 space-y-2 overflow-y-auto">
                {issuesByStatus[status]?.map((issue) => (
                  <Card
                    key={issue.id}
                    draggable
                    onDragStart={(e) => handleDragStart(e, issue)}
                    onClick={() => setSelectedIssue(issue)}
                    className="cursor-pointer hover:shadow-md transition-shadow"
                  >
                    <CardHeader className="p-3">
                      <div className="flex items-start justify-between gap-2">
                        <CardTitle className="text-sm font-medium line-clamp-2">
                          {issue.title}
                        </CardTitle>
                        <Badge className={PRIORITY_COLORS[issue.priority]} variant="outline">
                          {issue.priority}
                        </Badge>
                      </div>
                    </CardHeader>
                    <CardContent className="p-3 pt-0">
                      <div className="flex items-center justify-between text-xs text-muted-foreground">
                        <span>{issue.issue_key}</span>
                        {issue.story_points && (
                          <Badge variant="secondary">{issue.story_points} pts</Badge>
                        )}
                      </div>
                      <div className="mt-1 text-xs text-muted-foreground">
                        Assignee: {issue.assignee_id || 'Unassigned'}
                      </div>
                      {issue.labels.length > 0 && (
                        <div className="flex flex-wrap gap-1 mt-2">
                          {issue.labels.map((label) => (
                            <Badge key={label} variant="outline" className="text-xs">
                              {label}
                            </Badge>
                          ))}
                        </div>
                      )}
                    </CardContent>
                  </Card>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>

      {selectedIssue && (
        <IssueDialog
          issue={selectedIssue}
          open={!!selectedIssue}
          onClose={() => setSelectedIssue(null)}
        />
      )}

      <CreateIssueDialog
        projectId={projectId}
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
      />
    </div>
  );
}
