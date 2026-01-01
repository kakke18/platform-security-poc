interface PageHeaderProps {
  title: string;
}

export function PageHeader({ title }: PageHeaderProps) {
  return (
    <div className="bg-white border-b border-gray-200 px-8 py-6 mb-6 shadow-sm">
      <h1 className="text-3xl font-bold text-gray-900">{title}</h1>
    </div>
  );
}
