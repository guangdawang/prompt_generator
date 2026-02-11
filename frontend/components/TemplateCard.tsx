'use client';

import { Template } from '@/lib/api';

interface TemplateCardProps {
  template: Template;
  onSelect: (template: Template) => void;
  onEdit?: (template: Template) => void;
}

export default function TemplateCard({ template, onSelect, onEdit }: TemplateCardProps) {
  return (
    <div
      className="border border-gray-200 rounded-lg p-4 cursor-pointer hover:shadow-md transition-shadow duration-200 hover:border-blue-300 bg-white"
      onClick={() => onSelect(template)}
    >
      <div className="flex items-start justify-between mb-2 gap-3">
        <h3 className="font-semibold text-lg text-gray-900">{template.name}</h3>
        <div className="flex items-center gap-2">
          <span className="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded">
            {template.category || '未分类'}
          </span>
          {onEdit && (
            <button
              type="button"
              onClick={(event) => {
                event.stopPropagation();
                onEdit(template);
              }}
              className="px-2 py-1 text-xs bg-blue-50 text-blue-700 rounded hover:bg-blue-100"
            >
              编辑
            </button>
          )}
        </div>
      </div>
      <p className="text-gray-600 text-sm mb-3 line-clamp-2">{template.description}</p>
      <div className="flex items-center justify-between text-xs text-gray-500">
        <span>使用次数: {template.usage_count}</span>
        {template.is_public && (
          <span className="px-2 py-1 bg-green-100 text-green-700 rounded">
            公开
          </span>
        )}
      </div>
    </div>
  );
}
