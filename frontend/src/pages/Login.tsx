import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authService } from '@/lib/auth';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useToast } from '@/hooks/use-toast';
import { Zap } from 'lucide-react';

export default function Login() {
  const navigate = useNavigate();
  const { toast } = useToast();
  const [isLogin, setIsLogin] = useState(true);
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    display_name: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      if (isLogin) {
        await authService.login({
          email: formData.email,
          password: formData.password,
        });
        toast({
          title: 'Welcome back!',
          description: 'You have successfully logged in.',
        });
      } else {
        await authService.register({
          email: formData.email,
          password: formData.password,
          display_name: formData.display_name,
        });
        toast({
          title: 'Account created!',
          description: 'Your account has been created successfully.',
        });
      }
      navigate('/projects');
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.message || 'An error occurred. Please try again.',
        variant: 'destructive',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  return (
    <div className="h-screen overflow-hidden bg-gradient-to-br from-slate-50 via-white to-blue-50 text-slate-900">
      <div className="relative h-full">
        <div className="pointer-events-none absolute -left-16 top-8 h-64 w-64 rounded-full bg-blue-200/50 blur-3xl" />
        <div className="pointer-events-none absolute -right-16 bottom-8 h-64 w-64 rounded-full bg-cyan-200/50 blur-3xl" />

        <div className="mx-auto grid h-full max-w-7xl grid-cols-1 px-4 py-4 lg:grid-cols-2 lg:px-8 lg:py-8">
          <section className="hidden h-full flex-col justify-center pr-12 lg:flex">
            <div className="inline-flex w-fit items-center gap-2 rounded-full border border-blue-200 bg-blue-50 px-4 py-2 text-sm font-medium text-blue-700">
              <Zap className="h-4 w-4" />
              Streamline your workflow
            </div>

            <h1 className="mt-6 text-5xl font-bold leading-tight text-slate-900">
              Project Management
              <span className="mt-2 block text-blue-600">Made Simple</span>
            </h1>

            <p className="mt-6 max-w-xl text-lg leading-relaxed text-slate-600">
              Plan sprints, track issues, and collaborate in real time from one focused workspace.
            </p>

            <div className="relative mt-8 max-w-xl">
              <div className="relative z-10 rounded-2xl border border-slate-200 bg-white p-6 shadow-2xl transition-transform duration-300 hover:rotate-0 lg:rotate-2">
                <div className="space-y-4">
                  <div className="flex items-center justify-between border-b border-slate-100 pb-4">
                    <div className="flex items-center gap-2">
                      <div className="h-3 w-3 rounded-full bg-red-400" />
                      <div className="h-3 w-3 rounded-full bg-amber-400" />
                      <div className="h-3 w-3 rounded-full bg-emerald-400" />
                    </div>
                    <span className="text-xs text-slate-500">Project Board</span>
                  </div>

                  <div className="grid grid-cols-3 gap-3">
                    {['To Do', 'In Progress', 'Done'].map((status, idx) => (
                      <div key={status} className="space-y-2">
                        <div className="text-xs font-semibold text-slate-500">{status}</div>
                        {[1, 2].map((item) => (
                          <div
                            key={item}
                            className="rounded bg-slate-100 p-3 text-xs animate-pulse"
                            style={{ animationDelay: `${idx * 200 + item * 100}ms` }}
                          >
                            <div className="mb-2 h-2 rounded bg-slate-300" />
                            <div className="h-2 w-2/3 rounded bg-slate-200" />
                          </div>
                        ))}
                      </div>
                    ))}
                  </div>
                </div>
              </div>
              <div className="pointer-events-none absolute -left-6 -top-6 h-24 w-24 rounded-full bg-blue-100 blur-2xl" />
              <div className="pointer-events-none absolute -bottom-8 -right-8 h-28 w-28 rounded-full bg-cyan-100 blur-2xl" />
            </div>
          </section>

          <section className="flex h-full items-center justify-center">
            <Card className="w-full max-w-md border-slate-200 bg-white/95 shadow-2xl backdrop-blur">
              <CardHeader className="space-y-2">
                <CardTitle className="text-center text-2xl font-bold text-slate-900">
                  {isLogin ? 'Welcome back' : 'Create an account'}
                </CardTitle>
                <CardDescription className="text-center text-slate-600">
                  {isLogin
                    ? 'Enter your credentials to continue'
                    : 'Enter your details to get started'}
                </CardDescription>
              </CardHeader>

              <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                  {isLogin && (
                    <button
                      type="button"
                      onClick={() =>
                        setFormData({
                          email: 'user@gmail.com',
                          password: 'password1',
                          display_name: '',
                        })
                      }
                      className="w-full text-sm text-blue-600 hover:text-blue-700 underline"
                    >
                      Use demo account
                    </button>
                  )}
                  <div className="space-y-2">
                    <Label htmlFor="email" className="text-slate-700">
                      Email
                    </Label>
                    <Input
                      id="email"
                      name="email"
                      type="email"
                      placeholder="name@example.com"
                      value={formData.email}
                      onChange={handleChange}
                      required
                      disabled={isLoading}
                    />
                  </div>

                  {!isLogin && (
                    <div className="space-y-2">
                      <Label htmlFor="display_name" className="text-slate-700">
                        Display Name
                      </Label>
                      <Input
                        id="display_name"
                        name="display_name"
                        type="text"
                        placeholder="John Doe"
                        value={formData.display_name}
                        onChange={handleChange}
                        required
                        disabled={isLoading}
                      />
                    </div>
                  )}

                  <div className="space-y-2">
                    <Label htmlFor="password" className="text-slate-700">
                      Password
                    </Label>
                    <Input
                      id="password"
                      name="password"
                      type="password"
                      placeholder="••••••••"
                      value={formData.password}
                      onChange={handleChange}
                      required
                      minLength={8}
                      disabled={isLoading}
                    />
                    {!isLogin && (
                      <p className="text-xs text-slate-500">
                        Password must be at least 8 characters long
                      </p>
                    )}
                  </div>

                  <Button type="submit" className="w-full" disabled={isLoading}>
                    {isLoading ? 'Loading...' : isLogin ? 'Sign In' : 'Sign Up'}
                  </Button>
                </form>

                <div className="mt-4 text-center text-sm text-slate-600">
                  <button
                    type="button"
                    onClick={() => setIsLogin(!isLogin)}
                    className="text-blue-600 hover:underline"
                    disabled={isLoading}
                  >
                    {isLogin
                      ? "Don't have an account? Sign up"
                      : 'Already have an account? Sign in'}
                  </button>
                </div>
              </CardContent>
            </Card>
          </section>
        </div>
      </div>
    </div>
  );
}