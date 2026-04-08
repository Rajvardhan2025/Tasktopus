import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ArrowRight, CheckCircle, AlertCircle } from 'lucide-react';
import api from '@/lib/api';

interface Workflow {
  id: string;
  name: string;
  project_id: string;
  statuses: string[];
  transitions: Transition[];
}

interface Transition {
  from: string;
  to: string;
  conditions: Condition[];
  actions: Action[];
}

interface Condition {
  field: string;
  operator: string;
  value: string;
}

interface Action {
  type: string;
  params: Record<string, any>;
}

interface WorkflowManagementProps {
  projectId: string;
}

const STATUS_LABELS: Record<string, string> = {
  to_do: 'To Do',
  in_progress: 'In Progress',
  in_review: 'In Review',
  done: 'Done',
};

const STATUS_COLORS: Record<string, string> = {
  to_do: 'bg-gray-100 text-gray-800 border-gray-300',
  in_progress: 'bg-blue-100 text-blue-800 border-blue-300',
  in_review: 'bg-yellow-100 text-yellow-800 border-yellow-300',
  done: 'bg-green-100 text-green-800 border-green-300',
};

export function WorkflowManagement({ projectId }: WorkflowManagementProps) {
  const { data: workflowData, isLoading } = useQuery({
    queryKey: ['workflow', projectId],
    queryFn: async () => {
      const response = await api.get<{ data: Workflow }>(`/projects/${projectId}/workflow`);
      return response.data;
    },
  });

  const workflow = workflowData?.data;

  if (isLoading) {
    return <div className="flex items-center justify-center h-64">Loading workflow...</div>;
  }

  if (!workflow) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <p className="text-muted-foreground mb-4">No workflow configured for this project</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold mb-2">Workflow Configuration</h2>
        <p className="text-muted-foreground">{workflow.name}</p>
      </div>

      {/* Status Flow Visualization */}
      <Card>
        <CardHeader>
          <CardTitle>Status Flow</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-2 flex-wrap">
            {workflow.statuses.map((status, index) => (
              <div key={status} className="flex items-center gap-2">
                <Badge
                  variant="outline"
                  className={`${STATUS_COLORS[status]} px-4 py-2 text-sm font-medium`}
                >
                  {STATUS_LABELS[status] || status}
                </Badge>
                {index < workflow.statuses.length - 1 && (
                  <ArrowRight className="w-4 h-4 text-muted-foreground" />
                )}
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Transitions */}
      <div>
        <h3 className="text-lg font-semibold mb-3">Allowed Transitions</h3>
        <div className="space-y-3">
          {workflow.transitions.map((transition, index) => (
            <Card key={index}>
              <CardContent className="pt-6">
                <div className="space-y-4">
                  {/* Transition Path */}
                  <div className="flex items-center gap-3">
                    <Badge
                      variant="outline"
                      className={`${STATUS_COLORS[transition.from]} px-3 py-1`}
                    >
                      {STATUS_LABELS[transition.from] || transition.from}
                    </Badge>
                    <ArrowRight className="w-5 h-5 text-muted-foreground" />
                    <Badge
                      variant="outline"
                      className={`${STATUS_COLORS[transition.to]} px-3 py-1`}
                    >
                      {STATUS_LABELS[transition.to] || transition.to}
                    </Badge>
                  </div>

                  {/* Conditions */}
                  {transition.conditions && transition.conditions.length > 0 && (
                    <div className="pl-4 border-l-2 border-orange-300">
                      <div className="flex items-center gap-2 mb-2">
                        <AlertCircle className="w-4 h-4 text-orange-600" />
                        <span className="text-sm font-semibold text-orange-600">
                          Conditions Required
                        </span>
                      </div>
                      <div className="space-y-1">
                        {transition.conditions.map((condition, idx) => (
                          <div key={idx} className="text-sm text-muted-foreground">
                            • {formatCondition(condition)}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Actions */}
                  {transition.actions && transition.actions.length > 0 && (
                    <div className="pl-4 border-l-2 border-blue-300">
                      <div className="flex items-center gap-2 mb-2">
                        <CheckCircle className="w-4 h-4 text-blue-600" />
                        <span className="text-sm font-semibold text-blue-600">
                          Automatic Actions
                        </span>
                      </div>
                      <div className="space-y-1">
                        {transition.actions.map((action, idx) => (
                          <div key={idx} className="text-sm text-muted-foreground">
                            • {formatAction(action)}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>

      {/* Workflow Rules Summary */}
      <Card>
        <CardHeader>
          <CardTitle>Workflow Rules Summary</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3 text-sm">
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-600 mt-0.5" />
              <div>
                <div className="font-medium">Status Columns</div>
                <div className="text-muted-foreground">
                  {workflow.statuses.length} configurable status columns
                </div>
              </div>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-600 mt-0.5" />
              <div>
                <div className="font-medium">Transition Rules</div>
                <div className="text-muted-foreground">
                  {workflow.transitions.length} defined transitions with validation
                </div>
              </div>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-600 mt-0.5" />
              <div>
                <div className="font-medium">Validation Hooks</div>
                <div className="text-muted-foreground">
                  Conditions prevent invalid transitions (e.g., missing required fields)
                </div>
              </div>
            </div>
            <div className="flex items-start gap-2">
              <CheckCircle className="w-4 h-4 text-green-600 mt-0.5" />
              <div>
                <div className="font-medium">Automatic Actions</div>
                <div className="text-muted-foreground">
                  Actions execute automatically on status changes (e.g., notifications)
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function formatCondition(condition: Condition): string {
  const fieldLabels: Record<string, string> = {
    assignee_id: 'Assignee',
    story_points: 'Story Points',
    description: 'Description',
  };

  const operatorLabels: Record<string, string> = {
    not_empty: 'must be set',
    equals: 'must equal',
    greater_than: 'must be greater than',
  };

  const field = fieldLabels[condition.field] || condition.field;
  const operator = operatorLabels[condition.operator] || condition.operator;
  const value = condition.value ? ` ${condition.value}` : '';

  return `${field} ${operator}${value}`;
}

function formatAction(action: Action): string {
  const actionLabels: Record<string, string> = {
    notify: 'Send notification',
    assign_reviewer: 'Assign reviewer',
    set_field: 'Set field value',
  };

  const label = actionLabels[action.type] || action.type;
  const params = action.params?.message ? `: "${action.params.message}"` : '';

  return `${label}${params}`;
}
