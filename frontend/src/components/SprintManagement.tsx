import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Sprint, Issue, sprintsApi, issuesApi } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Plus, Play, CheckCircle, Calendar, Target } from 'lucide-react';
import { format } from 'date-fns';

interface SprintManagementProps {
  projectId: string;
}

export function SprintManagement({ projectId }: SprintManagementProps) {
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [completeDialogOpen, setCompleteDialogOpen] = useState(false);
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
        start_date: startDate,
        end_date: endDate,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sprints', projectId] });
      handleCloseCreate();
    },
  });

  const startSprintMutation = useMutation({
    mutationFn: (sprintId: string) => sprintsApi.start(sprintId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sprints', projectId] });
    },
  });

  const completeSprintMutation = useMutation({
    mutationFn: ({ sprintId, carryOver }: { sprintId: string; carryOver: string[] }) =>
      sprintsApi.complete(sprintId, carryOver),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sprints', projectId] });
      queryClient.invalidateQueries({ queryKey: ['issues', projectId] });
      handleCloseComplete();
    },
  });

  const sprints = sprintsData?.data?.data || [];
  const issues = issuesData?.data?.data?.issues || [];
  
  const activeSprint = sprints.find((s) => s.status === 'active');
  const backlogIssues = issues.filter((i) => !i.sprint_id);

  const handleCloseCreate = () => {
    setCreateDialogOpen(false);
    setSprintName('');
    setSprintGoal('');
    setStartDate('');
    setEndDate('');
  };

  const handleCloseComplete = () => {
    setCompleteDialogOpen(false);
    setSelectedSprint(null);
    setCarryOverIssues([]);
  };

  const handleCompleteSprint = (sprint: Sprint) => {
    setSelectedSprint(sprint);
    setCompleteDialogOpen(true);
  };

  const incompleteIssues = selectedSprint
    ? issues.filter((i) => i.sprint_id === selectedSprint.id && i.status !== 'done')
    : [];

  const toggleCarryOver = (issueId: string) => {
    setCarryOverIssues((prev) =>
      prev.includes(issueId) ? prev.filter((id) => id !== issueId) : [...prev, issueId]
    );
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
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="flex items-center gap-2">
                  <Play className="w-5 h-5 text-green-600" />
                  {activeSprint.name}
                  <Badge variant="default">Active</Badge>
                </CardTitle>
                <p className="text-sm text-muted-foreground mt-1">{activeSprint.goal}</p>
              </div>
              <Button onClick={() => handleCompleteSprint(activeSprint)}>
                <CheckCircle className="w-4 h-4 mr-2" />
                Complete Sprint
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-4 gap-4">
              <div>
                <div className="text-sm text-muted-foreground">Start Date</div>
                <div className="font-medium">
                  {format(new Date(activeSprint.start_date), 'MMM dd, yyyy')}
                </div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">End Date</div>
                <div className="font-medium">
                  {format(new Date(activeSprint.end_date), 'MMM dd, yyyy')}
                </div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">Completed Points</div>
                <div className="font-medium">
                  {activeSprint.completed_points} / {activeSprint.total_points}
                </div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">Velocity</div>
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
      )}

      {/* Planned Sprints */}
      <div>
        <h3 className="text-lg font-semibold mb-3">Planned Sprints</h3>
        <div className="space-y-3">
          {sprints
            .filter((s) => s.status === 'planned')
            .map((sprint) => (
              <Card key={sprint.id}>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle className="flex items-center gap-2">
                        <Calendar className="w-5 h-5" />
                        {sprint.name}
                        <Badge variant="secondary">Planned</Badge>
                      </CardTitle>
                      <p className="text-sm text-muted-foreground mt-1">{sprint.goal}</p>
                    </div>
                    <Button
                      variant="outline"
                      onClick={() => startSprintMutation.mutate(sprint.id)}
                      disabled={!!activeSprint}
                    >
                      <Play className="w-4 h-4 mr-2" />
                      Start Sprint
                    </Button>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="flex gap-4 text-sm">
                    <div>
                      <span className="text-muted-foreground">Start:</span>{' '}
                      {format(new Date(sprint.start_date), 'MMM dd, yyyy')}
                    </div>
                    <div>
                      <span className="text-muted-foreground">End:</span>{' '}
                      {format(new Date(sprint.end_date), 'MMM dd, yyyy')}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
        </div>
      </div>

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

      {/* Create Sprint Dialog */}
      <Dialog open={createDialogOpen} onOpenChange={handleCloseCreate}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create Sprint</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label htmlFor="name">Sprint Name</Label>
              <Input
                id="name"
                value={sprintName}
                onChange={(e) => setSprintName(e.target.value)}
                placeholder="Sprint 1"
              />
            </div>
            <div>
              <Label htmlFor="goal">Sprint Goal</Label>
              <Textarea
                id="goal"
                value={sprintGoal}
                onChange={(e) => setSprintGoal(e.target.value)}
                placeholder="What do you want to achieve?"
                rows={3}
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="startDate">Start Date</Label>
                <Input
                  id="startDate"
                  type="date"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                />
              </div>
              <div>
                <Label htmlFor="endDate">End Date</Label>
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
              disabled={!sprintName || !startDate || !endDate}
            >
              Create Sprint
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Complete Sprint Dialog */}
      <Dialog open={completeDialogOpen} onOpenChange={handleCloseComplete}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Complete Sprint</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <p className="text-sm text-muted-foreground">
              {incompleteIssues.length} incomplete issue(s) found. Select which issues to carry over
              to the next sprint.
            </p>
            <div className="space-y-2 max-h-96 overflow-y-auto">
              {incompleteIssues.map((issue) => (
                <div
                  key={issue.id}
                  className="flex items-center gap-3 p-3 border rounded hover:bg-muted/50"
                >
                  <input
                    type="checkbox"
                    checked={carryOverIssues.includes(issue.id)}
                    onChange={() => toggleCarryOver(issue.id)}
                    className="w-4 h-4"
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
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={handleCloseComplete}>
              Cancel
            </Button>
            <Button
              onClick={() =>
                selectedSprint &&
                completeSprintMutation.mutate({
                  sprintId: selectedSprint.id,
                  carryOver: carryOverIssues,
                })
              }
            >
              Complete Sprint
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
