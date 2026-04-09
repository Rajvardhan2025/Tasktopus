import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { 
  ArrowRight, 
  LayoutDashboard, 
  Zap, 
  Users, 
  GitBranch,
  CheckCircle2,
  TrendingUp,
  Bell
} from 'lucide-react';

export function Landing() {
  const navigate = useNavigate();

  const features = [
    {
      icon: LayoutDashboard,
      title: 'Kanban Boards',
      description: 'Visualize your workflow with intuitive drag-and-drop boards'
    },
    {
      icon: Zap,
      title: 'Sprint Management',
      description: 'Plan and track sprints with velocity metrics and burndown charts'
    },
    {
      icon: Users,
      title: 'Team Collaboration',
      description: 'Real-time updates and seamless team communication'
    },
    {
      icon: GitBranch,
      title: 'Custom Workflows',
      description: 'Design workflows that match your team\'s process'
    },
    {
      icon: Bell,
      title: 'Smart Notifications',
      description: 'Stay informed with intelligent activity tracking'
    },
    {
      icon: TrendingUp,
      title: 'Analytics',
      description: 'Track progress and optimize team performance'
    }
  ];

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="relative overflow-hidden bg-gradient-to-br from-primary/5 via-background to-secondary/5">
        <div className="container mx-auto px-4 py-20 md:py-32">
          <div className="grid lg:grid-cols-2 gap-12 items-center">
            <div className="space-y-8 animate-fade-in">
              <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 text-primary text-sm font-medium">
                <Zap className="w-4 h-4" />
                Streamline Your Workflow
              </div>
              
              <h1 className="text-5xl md:text-6xl font-bold leading-tight">
                Project Management
                <span className="block text-primary mt-2">Made Simple</span>
              </h1>
              
              <p className="text-xl text-muted-foreground leading-relaxed">
                Empower your team with intuitive project tracking, real-time collaboration, 
                and powerful insights. From planning to delivery, manage everything in one place.
              </p>
              
              <div className="flex flex-col sm:flex-row gap-4">
                <Button 
                  size="lg" 
                  onClick={() => navigate('/projects')}
                  className="text-lg px-8 py-6 group"
                >
                  Get Started
                  <ArrowRight className="ml-2 w-5 h-5 group-hover:translate-x-1 transition-transform" />
                </Button>
                <Button 
                  size="lg" 
                  variant="outline"
                  onClick={() => document.getElementById('features')?.scrollIntoView({ behavior: 'smooth' })}
                  className="text-lg px-8 py-6"
                >
                  Learn More
                </Button>
              </div>

              <div className="flex items-center gap-8 pt-4">
                <div className="flex items-center gap-2">
                  <CheckCircle2 className="w-5 h-5 text-green-500" />
                  <span className="text-sm text-muted-foreground">No credit card required</span>
                </div>
                <div className="flex items-center gap-2">
                  <CheckCircle2 className="w-5 h-5 text-green-500" />
                  <span className="text-sm text-muted-foreground">Free to start</span>
                </div>
              </div>
            </div>

            {/* Hero Image/Illustration */}
            <div className="relative lg:block hidden">
              <div className="relative z-10 bg-card border rounded-2xl shadow-2xl p-6 transform rotate-2 hover:rotate-0 transition-transform duration-300">
                <div className="space-y-4">
                  <div className="flex items-center justify-between pb-4 border-b">
                    <div className="flex items-center gap-2">
                      <div className="w-3 h-3 rounded-full bg-red-500"></div>
                      <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                      <div className="w-3 h-3 rounded-full bg-green-500"></div>
                    </div>
                    <span className="text-xs text-muted-foreground">Project Board</span>
                  </div>
                  
                  <div className="grid grid-cols-3 gap-4">
                    {['To Do', 'In Progress', 'Done'].map((status, idx) => (
                      <div key={status} className="space-y-2">
                        <div className="text-xs font-semibold text-muted-foreground">{status}</div>
                        {[1, 2].map((item) => (
                          <div 
                            key={item}
                            className="bg-muted/50 rounded p-3 text-xs animate-pulse"
                            style={{ animationDelay: `${idx * 200 + item * 100}ms` }}
                          >
                            <div className="h-2 bg-muted-foreground/20 rounded mb-2"></div>
                            <div className="h-2 bg-muted-foreground/10 rounded w-2/3"></div>
                          </div>
                        ))}
                      </div>
                    ))}
                  </div>
                </div>
              </div>
              
              {/* Decorative elements */}
              <div className="absolute -top-4 -right-4 w-72 h-72 bg-primary/10 rounded-full blur-3xl"></div>
              <div className="absolute -bottom-4 -left-4 w-72 h-72 bg-secondary/10 rounded-full blur-3xl"></div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-20 md:py-32">
        <div className="container mx-auto px-4">
          <div className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-4xl md:text-5xl font-bold mb-4">
              Everything You Need
            </h2>
            <p className="text-xl text-muted-foreground">
              Powerful features designed to help your team work smarter, not harder
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <Card 
                key={feature.title}
                className="group hover:shadow-lg transition-all duration-300 hover:-translate-y-1 cursor-pointer"
                style={{ animationDelay: `${index * 100}ms` }}
              >
                <CardContent className="p-6">
                  <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4 group-hover:bg-primary/20 transition-colors">
                    <feature.icon className="w-6 h-6 text-primary" />
                  </div>
                  <h3 className="text-xl font-semibold mb-2">{feature.title}</h3>
                  <p className="text-muted-foreground">{feature.description}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 md:py-32 bg-gradient-to-br from-primary/10 via-background to-secondary/10">
        <div className="container mx-auto px-4">
          <Card className="max-w-4xl mx-auto border-2">
            <CardContent className="p-12 text-center">
              <h2 className="text-4xl md:text-5xl font-bold mb-4">
                Ready to Get Started?
              </h2>
              <p className="text-xl text-muted-foreground mb-8">
                Join teams who are already managing their projects more efficiently
              </p>
              <Button 
                size="lg"
                onClick={() => navigate('/projects')}
                className="text-lg px-8 py-6 group"
              >
                Start Your First Project
                <ArrowRight className="ml-2 w-5 h-5 group-hover:translate-x-1 transition-transform" />
              </Button>
            </CardContent>
          </Card>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t py-8">
        <div className="container mx-auto px-4">
          <div className="flex flex-col md:flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-2">
              <LayoutDashboard className="w-5 h-5" />
              <span className="font-semibold">Project Management</span>
            </div>
            <p className="text-sm text-muted-foreground">
              © 2024 Project Management. Built with ❤️ for productive teams.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}
