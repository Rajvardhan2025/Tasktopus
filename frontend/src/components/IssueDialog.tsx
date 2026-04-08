import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Issue, commentsApi, activityApi, issuesApi } from '@/lib/api';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Textarea } from '@/components/ui/textarea';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card } from '@/components/ui/card';
import { formatDistanceToNow } from 'date-fns';
import { MessageSquare, Activity, Edit2, Save, X } from 'lucide-react';

interface IssueDialogProps {
  issue: Issue;
  open: boolean;
  onClose: () => void;
}

export function IssueDialog({ issue, open, onClose }: IssueDialogProps) {
  const [activeTab, setActiveTab] = useState<'comments' | 'activity'>('comments');
  const [commentText, setCommentText] = useState('');
  const [isEditing, setIsEditing] = useState(false);
  const [editedTitle, setEditedTitle] = useState(issue.title);
  const [editedDescription, setEditedDescription] = useState(issue.description);
  const [editedStoryPoints, setEditedStoryPoints] = useState(issue.story_points?.toString() || '');
  const [editedPriority, setEditedPriority] = useState(issue.priority);
  const queryClient = useQueryClient();

  const { data: commentsData } = useQuery({
    queryKey: ['comments', issue.id],
    queryFn: () => commentsApi.list(issue.id),
    enabled: open && activeTab === 'comments',
  });

  const { data: activityData } = useQuery({
    queryKey: ['activity', issue.id],
    queryFn: () => activityApi.issue(issue.id),
    enabled: open && activeTab === 'activity',
  });

  const addCommentMutation = useMutation({
    mutationFn: () => commentsApi.create(issue.id, commentText),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', issue.id] });
      setCommentText('');
    },
  });

  const updateIssueMutation = useMutation({
    mutationFn: () =>
      issuesApi.update(issue.id, {
        title: editedTitle,
        description: editedDescription,
        story_points: editedStoryPoints ? parseInt(editedStoryPoints) : undefined,
        priority: editedPriority,
        version: issue.version,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['issues', issue.project_id] });
      setIsEditing(false);
    },
  });

  const handleSave = () => {
    updateIssueMutation.mutate();
  };

  const handleCancelEdit = () => {
    setEditedTitle(issue.title);
    setEditedDescription(issue.description);
    setEditedStoryPoints(issue.story_points?.toString() || '');
    setEditedPriority(issue.priority);
    setIsEditing(false);
  };

  const comments = commentsData?.data?.data || [];
  const activities = activityData?.data?.data || [];

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <div className="flex items-start justify-between">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-2">
                <Badge variant="outline">{issue.issue_key}</Badge>
                <Badge>{issue.type}</Badge>
                <Badge variant="secondary">{issue.priority}</Badge>
              </div>
              {isEditing ? (
                <Input
                  value={editedTitle}
                  onChange={(e) => setEditedTitle(e.target.value)}
                  className="text-2xl font-bold"
                />
              ) : (
                <DialogTitle className="text-2xl">{issue.title}</DialogTitle>
              )}
            </div>
            <div className="flex gap-2">
              {isEditing ? (
                <>
                  <Button size="sm" onClick={handleSave} disabled={updateIssueMutation.isPending}>
                    <Save className="w-4 h-4 mr-2" />
                    Save
                  </Button>
                  <Button size="sm" variant="outline" onClick={handleCancelEdit}>
                    <X className="w-4 h-4 mr-2" />
                    Cancel
                  </Button>
                </>
              ) : (
                <Button size="sm" variant="outline" onClick={() => setIsEditing(true)}>
                  <Edit2 className="w-4 h-4 mr-2" />
                  Edit
                </Button>
              )}
            </div>
          </div>
        </DialogHeader>

        <div className="space-y-6">
          {/* Description */}
          <div>
            <h3 className="font-semibold mb-2">Description</h3>
            {isEditing ? (
              <Textarea
                value={editedDescription}
                onChange={(e) => setEditedDescription(e.target.value)}
                rows={4}
              />
            ) : (
              <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                {issue.description || 'No description provided'}
              </p>
            )}
          </div>

          {/* Details */}
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span className="font-semibold">Status:</span> {issue.status}
            </div>
            <div>
              <Label className="font-semibold">Priority:</Label>
              {isEditing ? (
                <select
                  value={editedPriority}
                  onChange={(e) => setEditedPriority(e.target.value as any)}
                  className="ml-2 border rounded px-2 py-1"
                >
                  <option value="lowest">Lowest</option>
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                  <option value="highest">Highest</option>
                </select>
              ) : (
                <span className="ml-2">{issue.priority}</span>
              )}
            </div>
            <div>
              <Label className="font-semibold">Story Points:</Label>
              {isEditing ? (
                <Input
                  type="number"
                  value={editedStoryPoints}
                  onChange={(e) => setEditedStoryPoints(e.target.value)}
                  className="ml-2 w-20 inline-block"
                  min="0"
                />
              ) : (
                <span className="ml-2">{issue.story_points || 'Not set'}</span>
              )}
            </div>
            {issue.labels.length > 0 && (
              <div>
                <span className="font-semibold">Labels:</span>
                <div className="flex flex-wrap gap-1 mt-1">
                  {issue.labels.map((label) => (
                    <Badge key={label} variant="outline" className="text-xs">
                      {label}
                    </Badge>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Tabs */}
          <div className="border-t pt-4">
            <div className="flex gap-4 mb-4">
              <Button
                variant={activeTab === 'comments' ? 'default' : 'ghost'}
                size="sm"
                onClick={() => setActiveTab('comments')}
              >
                <MessageSquare className="w-4 h-4 mr-2" />
                Comments ({comments.length})
              </Button>
              <Button
                variant={activeTab === 'activity' ? 'default' : 'ghost'}
                size="sm"
                onClick={() => setActiveTab('activity')}
              >
                <Activity className="w-4 h-4 mr-2" />
                Activity ({activities.length})
              </Button>
            </div>

            {activeTab === 'comments' && (
              <div className="space-y-4">
                <div>
                  <Textarea
                    placeholder="Add a comment... (use @username to mention)"
                    value={commentText}
                    onChange={(e) => setCommentText(e.target.value)}
                    rows={3}
                  />
                  <Button
                    className="mt-2"
                    onClick={() => addCommentMutation.mutate()}
                    disabled={!commentText.trim() || addCommentMutation.isPending}
                  >
                    Add Comment
                  </Button>
                </div>

                <div className="space-y-3">
                  {comments.map((comment) => (
                    <Card key={comment.id} className="p-4">
                      <div className="flex items-start justify-between mb-2">
                        <span className="font-semibold text-sm">{comment.user_id}</span>
                        <span className="text-xs text-muted-foreground">
                          {formatDistanceToNow(new Date(comment.created_at), { addSuffix: true })}
                        </span>
                      </div>
                      <p className="text-sm whitespace-pre-wrap">{comment.content}</p>
                    </Card>
                  ))}
                </div>
              </div>
            )}

            {activeTab === 'activity' && (
              <div className="space-y-3">
                {activities.map((activity) => (
                  <Card key={activity.id} className="p-4">
                    <div className="flex items-start justify-between mb-2">
                      <span className="font-semibold text-sm">{activity.user_id}</span>
                      <span className="text-xs text-muted-foreground">
                        {formatDistanceToNow(new Date(activity.timestamp), { addSuffix: true })}
                      </span>
                    </div>
                    <p className="text-sm">
                      <span className="font-medium">{activity.action}</span>
                      {activity.changes && (
                        <span className="text-muted-foreground ml-2">
                          {JSON.stringify(activity.changes)}
                        </span>
                      )}
                    </p>
                  </Card>
                ))}
              </div>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
