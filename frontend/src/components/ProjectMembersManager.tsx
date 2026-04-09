import { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { projectsApi, usersApi } from '@/lib/api';
import type { User } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useToast } from '@/hooks/use-toast';
import { Users, UserPlus, X } from 'lucide-react';

interface ProjectMembersManagerProps {
    projectId: string;
}

export function ProjectMembersManager({ projectId }: ProjectMembersManagerProps) {
    const [open, setOpen] = useState(false);
    const [selectedUserId, setSelectedUserId] = useState('');
    const [newName, setNewName] = useState('');
    const [newEmail, setNewEmail] = useState('');
    const { toast } = useToast();
    const queryClient = useQueryClient();

    const { data: membersData } = useQuery({
        queryKey: ['project-members', projectId],
        queryFn: () => projectsApi.members(projectId),
        enabled: !!projectId,
    });

    const { data: usersData } = useQuery({
        queryKey: ['users'],
        queryFn: usersApi.list,
        enabled: open,
    });

    const members = membersData?.data?.data || [];
    const users = usersData?.data?.data || [];

    const availableUsers = useMemo(
        () => users.filter((user) => !members.some((member) => member.id === user.id)),
        [users, members]
    );

    const addMemberMutation = useMutation({
        mutationFn: (userId: string) => projectsApi.addMember(projectId, userId),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['project-members', projectId] });
            setSelectedUserId('');
            toast({ title: 'Member added', variant: 'success' });
        },
        onError: (error: any) => {
            toast({
                title: 'Failed to add member',
                description: error.response?.data?.error?.message || 'An error occurred',
                variant: 'destructive',
            });
        },
    });

    const removeMemberMutation = useMutation({
        mutationFn: (userId: string) => projectsApi.removeMember(projectId, userId),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['project-members', projectId] });
            toast({ title: 'Member removed', variant: 'success' });
        },
        onError: (error: any) => {
            toast({
                title: 'Failed to remove member',
                description: error.response?.data?.error?.message || 'An error occurred',
                variant: 'destructive',
            });
        },
    });

    const createUserMutation = useMutation({
        mutationFn: () => usersApi.create({ email: newEmail, display_name: newName }),
        onSuccess: (response) => {
            queryClient.invalidateQueries({ queryKey: ['users'] });
            setSelectedUserId(response.data.data.id);
            setNewName('');
            setNewEmail('');
            toast({ title: 'User created', variant: 'success' });
        },
        onError: (error: any) => {
            toast({
                title: 'Failed to create user',
                description: error.response?.data?.error?.message || 'An error occurred',
                variant: 'destructive',
            });
        },
    });

    return (
        <>
            <Button variant="outline" onClick={() => setOpen(true)} className="gap-2">
                <Users className="h-4 w-4" />
                Members ({members.length})
            </Button>

            <Dialog open={open} onOpenChange={setOpen}>
                <DialogContent className="max-w-2xl">
                    <DialogHeader>
                        <DialogTitle>Project Members</DialogTitle>
                    </DialogHeader>

                    <div className="space-y-4">
                        <div>
                            <Label>Current Members</Label>
                            <div className="mt-2 space-y-2">
                                {members.length === 0 && (
                                    <div className="text-sm text-muted-foreground">No members added yet.</div>
                                )}
                                {members.map((member: User) => (
                                    <div key={member.id} className="flex items-center justify-between rounded-md border p-2">
                                        <div>
                                            <div className="font-medium text-sm">{member.display_name}</div>
                                            <div className="text-xs text-muted-foreground">{member.email}</div>
                                        </div>
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onClick={() => removeMemberMutation.mutate(member.id)}
                                        >
                                            <X className="h-4 w-4" />
                                        </Button>
                                    </div>
                                ))}
                            </div>
                        </div>

                        <div className="rounded-md border p-3">
                            <Label>Add Existing User</Label>
                            <div className="mt-2 flex items-center gap-2">
                                <select
                                    value={selectedUserId}
                                    onChange={(e) => setSelectedUserId(e.target.value)}
                                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                                >
                                    <option value="">Select user</option>
                                    {availableUsers.map((user: User) => (
                                        <option key={user.id} value={user.id}>
                                            {user.display_name} ({user.email})
                                        </option>
                                    ))}
                                </select>
                                <Button
                                    onClick={() => selectedUserId && addMemberMutation.mutate(selectedUserId)}
                                    disabled={!selectedUserId || addMemberMutation.isPending}
                                >
                                    <UserPlus className="h-4 w-4" />
                                </Button>
                            </div>
                        </div>

                        <div className="rounded-md border p-3">
                            <Label>Create and Add New User</Label>
                            <div className="mt-2 grid grid-cols-1 gap-2 md:grid-cols-2">
                                <Input placeholder="Display name" value={newName} onChange={(e) => setNewName(e.target.value)} />
                                <Input placeholder="Email" value={newEmail} onChange={(e) => setNewEmail(e.target.value)} />
                            </div>
                            <Button
                                className="mt-2"
                                variant="secondary"
                                onClick={() => createUserMutation.mutate()}
                                disabled={!newName || !newEmail || createUserMutation.isPending}
                            >
                                Create User
                            </Button>
                        </div>

                        <div className="flex flex-wrap gap-2">
                            {members.map((member: User) => (
                                <Badge key={member.id} variant="secondary">
                                    {member.display_name}
                                </Badge>
                            ))}
                        </div>
                    </div>

                    <DialogFooter>
                        <Button variant="outline" onClick={() => setOpen(false)}>
                            Close
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    );
}
