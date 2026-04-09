import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { sprintsApi, issuesApi } from '@/lib/api';
import type { Sprint } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Plus, Play, CheckCircle, Calendar, AlertCircle } from 'lucide-react';
import { format } from 'date-fns';
import { useToast } from '@/hooks/use-toast';

interface SprintManagementProps {
  projectId: string;
}

export function SprintManagement({ projectId }: SprintManagementProps) {
  const { toast } = useToast();
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [closeDialogOpen, setCloseDialogOpen] = useState(false);
  const [selectedSprint, setSelectedSprint] = useState<Sprint | null>(null);
  const [sprintName, setSprintName] = useState('');
  const [sprintGoal, setSprintGoal] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [carryOverIssues, setCarryOverIssues] = useState<string[]>([]);

  const queryClient = useQueryClient();

  const { data: sprintsData } = useQuery({
    queryKey: ['sprints', projectId],
    queryFn: () => sprintsApi.list(projectId),
  });

  const { data: issuesData } = useQuery({
    queryKey: ['issues', projectId],
    queryFn: () => issuesApi.list(projectId),
  });

  const createSprintMutation = useMutation({
    mutationFn: () =>
      sprintsApi.create({
        project_id: projectId,
        name: sprintName,
        goal: sprintGoal,
        start_date: startDate ? new Date(startDate).toISOString() : undefined,
        end_date: endDate ? new Date(endDate).toISOString() : undefined,
      }),
    onSuccess: () => {
      toast({
        title: 'Sprint Created',
        description: `Sprint "${sprintName}" created successfully`,
        variant: 'success',
      });
      queryClient.invalidateQueries({ queryKey: ['sprints', projectId] });
      queryClient.invalidateQueries({ queryKey: ['issues', projectId] });
      handleCloseCreate();
    },
    onError: (error: any) => {
      toast({
        title: 'Failed to Create Sprint',
        description: error.response?.data?.message || 'An error occurred',
        variant: 'destructive',
      });
    },
  });

  const startSprintMutation = useMutation({
    mutationFn: (sprintId: string) => sprintsApi.start(sprintId),
    onSuccess: (_data: any, sprintId: string) => {
      const sprint = sprints.find(s => s.id === sprintId);
      toast({
        title: 'Sprint Started',
        description: `Sprint "${sprint?.name}" is now active`,
        variant: 'success',
      });
      queryClient.invalidateQueries({ queryKey: ['sprints', projectId] });
      queryClient.invalidateQueries({ queryKey: ['issues', projectId] });
    },
    onError: (error: any) => {
      toast({
        title: 'Cannot Start Sprint',
        description: error.response?.data?.message || 'An error occurred while starting the sprint',
        variant: 'destructive',
      });
    },
  });

  const closeSprintMutation = useMutation({
    mutationFn: ({ sprintId, carryOver }: { sprintId: string; carryOver: string[] }) =>
      sprintsApi.complete(sprintId, carryOver),
    onSuccess: () => {
      toast({
        title: 'Sprint Closed',
        description: 'Sprint has been closed successfully',
        variant: 'success',
      });
      queryClient.invalidateQueries({ queryKey: ['sprints', projectId] });
      queryClient.invalidateQueries({ queryKey: ['issues', projectId] });
      handleCloseClose();
    },
    onError: (error: any) => {
      toast({
        title: 'Failed to Close Sprint',
        description: error.response?.data?.message || 'An error occurred while closing the sprint',
        variant: 'destructive',
      });
    },
  });

  const sprints = sprintsData?.data?.data || [];
  const issues = issuesData?.data?.data?.issues || [];

  const activeSprint = sprints.find((s) => s.status === 'active');
  const futureSprints = sprints.filter((s) => s.status === 'future');
  const closedSprints = sprints.filter((s) => s.status === 'closed');
  const backlogIssues = issues.filter((i) => !i.sprint_id);

  const handleCloseCreate = () => {
    setCreateDialogOpen(false);
    setSprintName('');
    setSprintGoal('');
    setStartDate('');
    setEndDate('');
  };

  const handleCloseClose = () => {
    setCloseDialogOpen(false);
    setSelectedSprint(null);
    setCarryOverIssues([]);
  };

  const handleCloseSprint = (sprint: Sprint) => {
    setSelectedSprint(sprint);
    setCloseDialogOpen(true);
  };

  const incompleteIssues = selectedSprint
    ? issues.filter((i) => i.sprint_id === selectedSprint.id && i.status !== 'done')
    : [];

  const toggleCarryOver = (issueId: string) => {
    setCarryOverIssues((prev) =>
      prev.includes(issueId) ? prev.filter((id) => id !== issueId) : [...prev, issueId]
    );
  };

  const getSprintStatistics = (sprint: Sprint) => {
    const sprintIssues = issues.filter(i => i.sprint_id === sprint.id);
    const completedCount = sprintIssues.filter(i => i.status === 'done').length;
    const totalCount = sprintIssues.length;
    return { completedCount, totalCount };
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Sprint Management</h2>
        <Button onClick={() => setCreateDialogOpen(true)}>
          <Plus className="w-4 h-4 mr-2" />
          Create Sprint
        </Button>
      </div>

      {/* Active Sprint */}
      {activeSprint && (
        <>
          <Card className="border-green-200 bg-green-50">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="flex items-center gap-2">
                    <Play className="w-5 h-5 text-green-600" />
                    {activeSprint.name}
                    <Badge variant="default" className="bg-green-600">Active</Badge>
                  </CardTitle>
                  <p className="text-sm text-muted-foreground mt-1">{activeSprint.goal}</p>
                </div>
                <Button
                  variant="outline"
                  onClick={() => handleCloseSprint(activeSprint)}
                  className="border-green-200"
                >
                  <CheckCircle className="w-4 h-4 mr-2" />
                  Close Sprint
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-5 gap-4">
                <div>
                  <div className="text-sm text-muted-foreground">Start Date</div>
                  <div className="font-medium">
                    {activeSprint.start_date
                      ? format(new Date(activeSprint.start_date), 'MMM dd')
                      : 'Not set'}
                  </div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground">End Date</div>
                  <div className="font-medium">
                    {activeSprint.end_date
                      ? format(new Date(activeSprint.end_date), 'MMM dd')
                      : 'Not set'}
                  </div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground">Points</div>
                  <div className="font-medium">
                    {activeSprint.completed_points} / {activeSprint.total_points}
                  </div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground">Issues</div>
                  <div className="font-medium">
                    {getSprintStatistics(activeSprint).completedCount} / {getSprintStatistics(activeSprint).totalCount}
                  </div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground">Progress</div>
                  <div className="font-medium">
                    {activeSprint.total_points > 0
                      ? Math.round((activeSprint.completed_points / activeSprint.total_points) * 100)
                      : 0}
                    %
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </>
      )}

      {!activeSprint && futureSprints.length > 0 && (
        <Card className="border-yellow-200 bg-yellow-50 p-4">
          <div className="flex items-center gap-2">
            <AlertCircle className="w-5 h-5 text-yellow-600" />
            <p className="text-sm text-yellow-800">No active sprint. Start one to begin work.</p>
          </div>
        </Card>
      )}

      {/* Future Sprints */}
      {futureSprints.length > 0 && (
        <div>
          <h3 className="text-lg font-semibold mb-3">Upcoming Sprints ({futureSprints.length})</h3>
          <div className="space-y-3">
            {futureSprints.map((sprint) => (
              <Card key={sprint.id}>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle className="flex items-center gap-2">
                        <Calendar className="w-5 h-5" />
                        {sprint.name}
                        <Badge variant="secondary">Future</Badge>
                      </CardTitle>
                      <p className="text-sm text-muted-foreground mt-1">{sprint.goal}</p>
                    </div>
                    <Button
                      variant="outline"
                      onClick={() => startSprintMutation.mutate(sprint.id)}
                      disabled={startSprintMutation.isPending}
                    >
                      <Play className="w-4 h-4 mr-2" />
                      Start
                    </Button>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="flex gap-4 text-sm">
                    <div>
                      <span className="text-muted-foreground">Start:</span>{' '}
                      {sprint.start_date
                        ? format(new Date(sprint.start_date), 'MMM dd, yyyy')
                        : 'Not set'}
                    </div>
                    <div>
                      <span className="text-muted-foreground">End:</span>{' '}
                      {sprint.end_date
                        ? format(new Date(sprint.end_date), 'MMM dd, yyyy')
                        : 'Not set'}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}

      {/* Backlog */}
      <div>
        <h3 className="text-lg font-semibold mb-3">
          Backlog ({backlogIssues.length} issues)
        </h3>
        <Card>
          <CardContent className="pt-6">
            <div className="space-y-2">
              {backlogIssues.slice(0, 10).map((issue) => (
                <div
                  key={issue.id}
                  className="flex items-center justify-between p-2 hover:bg-muted/50 rounded"
                >
                  <div className="flex items-center gap-2">
                    <Badge variant="outline">{issue.issue_key}</Badge>
                    <span className="text-sm">{issue.title}</span>
                  </div>
                  {issue.story_points && (
                    <Badge variant="secondary">{issue.story_points} pts</Badge>
                  )}
                </div>
              ))}
              {backlogIssues.length > 10 && (
                <div className="text-sm text-muted-foreground text-center pt-2">
                  +{backlogIssues.length - 10} more issues
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Closed Sprints */}
      {closedSprints.length > 0 && (
        <div>
          <h3 className="text-lg font-semibold mb-3">Closed Sprints ({closedSprints.length})</h3>
          <div className="space-y-3">
            {closedSprints.map((sprint) => (
              <Card key={sprint.id} className="opacity-75">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <CheckCircle className="w-5 h-5 text-gray-400" />
                    {sprint.name}
                    <Badge variant="outline">Closed</Badge>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-4 gap-4 text-sm">
                    <div>
                      <span className="text-muted-foreground">Points:</span> {sprint.completed_points}/{sprint.total_points}
                    </div>
                    <div>
                      <span className="text-muted-foreground">Velocity:</span> {sprint.velocity}
                    </div>
                    <div>
                      <span className="text-muted-foreground">Issues:</span> {sprint.complete_issue_count}/{sprint.issue_count}
                    </div>
                    <div>
                      <span className="text-muted-foreground">Closed:</span> {sprint.end_date ? format(new Date(sprint.end_date), 'MMM dd') : 'N/A'}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}

      {/* Create Sprint Dialog */}
      <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create Sprint</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label htmlFor="name">Sprint Name *</Label>
              <Input
                id="name"
                value={sprintName}
                onChange={(e) => setSprintName(e.target.value)}
                placeholder="e.g., Sprint 1, Sprint 2"
              />
            </div>
            <div>
              <Label htmlFor="goal">Sprint Goal</Label>
              <Textarea
                id="goal"
                value={sprintGoal}
                onChange={(e) => setSprintGoal(e.target.value)}
                placeholder="What do you want to achieve in this sprint?"
                rows={3}
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="startDate">Start Date (optional)</Label>
                <Input
                  id="startDate"
                  type="date"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                />
              </div>
              <div>
                <Label htmlFor="endDate">End Date (optional)</Label>
                <Input
                  id="endDate"
                  type="date"
                  value={endDate}
                  onChange={(e) => setEndDate(e.target.value)}
                />
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={handleCloseCreate}>
              Cancel
            </Button>
            <Button
              onClick={() => createSprintMutation.mutate()}
              disabled={!sprintName || createSprintMutation.isPending}
            >
              {createSprintMutation.isPending ? 'Creating...' : 'Create Sprint'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Close Sprint Dialog */}
      <Dialog open={closeDialogOpen} onOpenChange={setCloseDialogOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Close Sprint</DialogTitle>
          </DialogHeader>
          {selectedSprint && (
            <div className="space-y-4">
              <div className="p-3 bg-blue-50 border border-blue-200 rounded">
                <p className="text-sm font-medium text-blue-900">
                  Closing: <strong>{selectedSprint.name}</strong>
                </p>
              </div>

              {incompleteIssues.length > 0 ? (
                <>
                  <p className="text-sm text-muted-foreground">
                    {incompleteIssues.length} incomplete issue(s) found. Select which issues to carry over to the next sprint.
                  </p>
                  <div className="space-y-2 max-h-96 overflow-y-auto border rounded">
                    {incompleteIssues.map((issue) => (
                      <div
                        key={issue.id}
                        className="flex items-center gap-3 p-3 border-b last:border-b-0 hover:bg-muted/30"
                      >
                        <input
                          type="checkbox"
                          checked={carryOverIssues.includes(issue.id)}
                          onChange={() => toggleCarryOver(issue.id)}
                          className="w-4 h-4 rounded border-gray-300"
                        />
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <Badge variant="outline">{issue.issue_key}</Badge>
                            <span className="text-sm font-medium">{issue.title}</span>
                          </div>
                          <div className="text-xs text-muted-foreground mt-1">
                            Status: {issue.status}
                          </div>
                        </div>
                        {issue.story_points && (
                          <Badge variant="secondary">{issue.story_points} pts</Badge>
                        )}
                      </div>
                    ))}
                  </div>
                </>
              ) : (
                <div className="p-3 bg-green-50 border border-green-200 rounded">
                  <p className="text-sm text-green-900">
                    ✓ All issues in this sprint are completed!
                  </p>
                </div>
              )}
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={handleCloseClose}>
              Cancel
            </Button>
            <Button
              onClick={() =>
                selectedSprint &&
                closeSprintMutation.mutate({
                  sprintId: selectedSprint.id,
                  carryOver: carryOverIssues,
                })
              }
              disabled={closeSprintMutation.isPending}
              variant="destructive"
            >
              {closeSprintMutation.isPending ? 'Closing...' : 'Close Sprint'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
