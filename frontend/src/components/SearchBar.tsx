import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { searchApi } from '@/lib/api';
import type { Issue } from '@/lib/api';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Card } from '@/components/ui/card';
import { Search, X } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface SearchBarProps {
  projectId?: string;
  onSelectIssue?: (issue: Issue) => void;
}

export function SearchBar({ projectId, onSelectIssue }: SearchBarProps) {
  const [query, setQuery] = useState('');
  const [isOpen, setIsOpen] = useState(false);

  const { data, isLoading } = useQuery({
    queryKey: ['search', query, projectId],
    queryFn: () => searchApi.search(query, projectId ? { project: projectId } : {}),
    enabled: query.length > 2,
  });

  const issues = data?.data?.data?.items || [];

  const handleSelect = (issue: Issue) => {
    onSelectIssue?.(issue);
    setQuery('');
    setIsOpen(false);
  };

  return (
    <div className="relative w-full max-w-md">
      <div className="relative">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
        <Input
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setIsOpen(true);
          }}
          onFocus={() => setIsOpen(true)}
          placeholder="Search issues... (e.g., status=in_progress priority>=high)"
          className="pl-10 pr-10"
        />
        {query && (
          <Button
            variant="ghost"
            size="sm"
            className="absolute right-1 top-1/2 transform -translate-y-1/2 h-7 w-7 p-0"
            onClick={() => {
              setQuery('');
              setIsOpen(false);
            }}
          >
            <X className="w-4 h-4" />
          </Button>
        )}
      </div>

      {isOpen && query.length > 2 && (
        <Card className="absolute top-full mt-2 w-full max-h-96 overflow-y-auto z-50 shadow-lg">
          {isLoading ? (
            <div className="p-4 text-center text-sm text-muted-foreground">Searching...</div>
          ) : issues.length === 0 ? (
            <div className="p-4 text-center text-sm text-muted-foreground">No issues found</div>
          ) : (
            <div className="p-2">
              {issues.map((issue) => (
                <div
                  key={issue.id}
                  onClick={() => handleSelect(issue)}
                  className="p-3 hover:bg-muted rounded cursor-pointer"
                >
                  <div className="flex items-start justify-between gap-2">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <Badge variant="outline" className="text-xs">
                          {issue.issue_key}
                        </Badge>
                        <Badge variant="secondary" className="text-xs">
                          {issue.type}
                        </Badge>
                        <Badge className="text-xs">{issue.priority}</Badge>
                      </div>
                      <div className="text-sm font-medium">{issue.title}</div>
                      <div className="text-xs text-muted-foreground mt-1">
                        Status: {issue.status}
                      </div>
                    </div>
                    {issue.story_points && (
                      <Badge variant="outline">{issue.story_points} pts</Badge>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>
      )}

      {/* Search Tips */}
      {isOpen && query.length <= 2 && query.length > 0 && (
        <Card className="absolute top-full mt-2 w-full p-4 z-50 shadow-lg">
          <div className="text-sm text-muted-foreground">
            <div className="font-semibold mb-2">Search Tips:</div>
            <ul className="space-y-1 text-xs">
              <li>• Type at least 3 characters to search</li>
              <li>• Use filters: status=in_progress</li>
              <li>• Priority threshold: priority&gt;=high</li>
              <li>• Filter by assignee: assignee=user123</li>
            </ul>
          </div>
        </Card>
      )}
    </div>
  );
}
