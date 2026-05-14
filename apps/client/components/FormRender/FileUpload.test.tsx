import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import FileUpload from './FileUpload'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useMutation } from '@tanstack/react-query'

// Mock useMutation
vi.mock('@tanstack/react-query', () => ({
  useMutation: vi.fn(),
}))

describe('FileUpload', () => {
  const mockMutate = vi.fn()
  const mockReset = vi.fn()
  const mockHandleResponse = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    ;(useMutation as any).mockReturnValue({
      mutate: mockMutate,
      reset: mockReset,
      isPending: false,
      isSuccess: false,
      isError: false,
      data: null,
      error: null,
    })
  })

  it('renders upload zone', () => {
    render(<FileUpload uploadPath="test-path" token="test-token" />)
    expect(screen.getByText(/drop your file here/i)).toBeDefined()
  })

  it('validates file type restrictions', async () => {
    const { container } = render(
      <FileUpload 
        uploadPath="test-path" 
        token="test-token" 
        accept="image/*" 
      />
    )
    
    const file = new File(['hello'], 'test.txt', { type: 'text/plain' })
    const input = container.querySelector('input[type="file"]') as HTMLInputElement
    
    fireEvent.change(input, { target: { files: [file] } })

    expect(await screen.findByText(/invalid file type/i)).toBeDefined()
    expect(mockMutate).not.toHaveBeenCalled()
  })

  it('allows valid file types', async () => {
    const { container } = render(
      <FileUpload 
        uploadPath="test-path" 
        token="test-token" 
        accept="image/*" 
      />
    )
    
    const file = new File(['hello'], 'test.png', { type: 'image/png' })
    const input = container.querySelector('input[type="file"]') as HTMLInputElement
    
    fireEvent.change(input, { target: { files: [file] } })

    expect(screen.queryByText(/invalid file type/i)).toBeNull()
    expect(screen.getByText('test.png')).toBeDefined()
  })

  it('triggers upload when manual button is clicked', async () => {
    const { container } = render(
      <FileUpload 
        uploadPath="test-path" 
        token="test-token" 
      />
    )
    
    const file = new File(['hello'], 'test.png', { type: 'image/png' })
    const input = container.querySelector('input[type="file"]') as HTMLInputElement
    
    fireEvent.change(input, { target: { files: [file] } })
    
    const uploadButton = screen.getByRole('button', { name: /upload file/i })
    fireEvent.click(uploadButton)

    expect(mockMutate).toHaveBeenCalledWith({ file, path: 'test-path' })
  })

  it('starts upload automatically when autoUpload is true', async () => {
    const { container } = render(
      <FileUpload 
        uploadPath="test-path" 
        token="test-token" 
        autoUpload={true}
      />
    )
    
    const file = new File(['hello'], 'test.png', { type: 'image/png' })
    const input = container.querySelector('input[type="file"]') as HTMLInputElement
    
    fireEvent.change(input, { target: { files: [file] } })

    await waitFor(() => {
      expect(mockMutate).toHaveBeenCalledWith({ file, path: 'test-path' })
    })
  })

  it('displays progress during upload', () => {
    ;(useMutation as any).mockReturnValue({
      mutate: mockMutate,
      reset: mockReset,
      isPending: true,
      isSuccess: false,
      isError: false,
    })
    
    // We need to trigger some progress state. 
    // Since progress is internal state updated by onProgress callback in mutationFn, 
    // we might need a more complex mock or just check if it renders when isPending is true.
    
    render(<FileUpload uploadPath="test-path" token="test-token" />)
    
    expect(screen.getByText(/uploading…/i)).toBeDefined()
  })

  it('calls handleResponse on success', async () => {
    const fileInfo = { fileName: 'test.png', filePath: '/path/test.png', fileSize: 100, fileType: 'image/png' }
    ;(useMutation as any).mockImplementation(({ onSuccess }) => {
      // Simulate success immediately for testing callback
      return {
        mutate: () => onSuccess({ data: { fileInfo } }),
        reset: mockReset,
        isPending: false,
        isSuccess: true,
        data: { data: { fileInfo } },
      }
    })

    const { container } = render(
      <FileUpload 
        uploadPath="test-path" 
        token="test-token" 
        handleResponse={mockHandleResponse}
        autoUpload={true}
      />
    )
    
    const file = new File(['hello'], 'test.png', { type: 'image/png' })
    const input = container.querySelector('input[type="file"]') as HTMLInputElement
    fireEvent.change(input, { target: { files: [file] } })

    expect(mockHandleResponse).toHaveBeenCalledWith(fileInfo)
  })
})
