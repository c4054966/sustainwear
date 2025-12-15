import { Button } from '@/components/ui/Button';

export default function LoginPage() {
    return (
        <div className="flex min-h-screen flex-col items-center justify-center bg-gray-100">
            <div className="w-full max-w-md rounded bg-white p-8 shadow">

                <h1 className="mb-6 text-2xl font-bold">Login</h1>
                <form className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700">
                            Email
                        </label>
                        <input
                            type="email"
                            className="mt-1 w-full rounded border border-gray-300 px-3 py-2"
                            placeholder="Enter your email"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700">
                            Password
                        </label>
                        <input
                            type="password"
                            className="mt-1 w-full rounded border border-gray-300 px-3 py-2"
                            placeholder="Enter your password"
                        />
                    </div>
                    <button
                        type="submit"
                        className="w-full rounded bg-blue-600 py-2 text-white font-medium hover:bg-blue-700"
                    >
                        Sign In
                    </button>
                    <Button variant="link" size="sm" className="w-full text-center">
                        Forgot Password?
                    </Button>
                </form>
            </div>
        </div>
    );
}