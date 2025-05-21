class Nessi < Formula
  desc "CLI-only Delta Lake quality and management tool"
  homepage "https://github.com/nessi-dev/nessi"
  version "1.0.0"
  license "Apache-2.0"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/nessi-dev/nessi/releases/download/v1.0.0/nessi_v1.0.0_darwin_arm64.tar.gz"
      sha256 "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
    else
      url "https://github.com/nessi-dev/nessi/releases/download/v1.0.0/nessi_v1.0.0_darwin_amd64.tar.gz"
      sha256 "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
    end
  elsif OS.linux?
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/nessi-dev/nessi/releases/download/v1.0.0/nessi_v1.0.0_linux_arm64.tar.gz"
      sha256 "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
    else
      url "https://github.com/nessi-dev/nessi/releases/download/v1.0.0/nessi_v1.0.0_linux_amd64.tar.gz"
      sha256 "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
    end
  end

  def install
    bin.install "nessi"
    # Generate and install shell completion scripts
    output = Utils.safe_popen_read("#{bin}/nessi", "completion", "bash")
    (bash_completion/"nessi").write output
    output = Utils.safe_popen_read("#{bin}/nessi", "completion", "zsh")
    (zsh_completion/"_nessi").write output
    output = Utils.safe_popen_read("#{bin}/nessi", "completion", "fish")
    (fish_completion/"nessi.fish").write output
  end

  test do
    assert_match "Nessi CLI v#{version}", shell_output("#{bin}/nessi version")
  end

  # Add a caveats section with helpful information
  def caveats
    <<~EOS
      Thank you for installing Nessi!
      
      To get started, run:
        nessi help
      
      For documentation, visit:
        https://github.com/nessi-dev/nessi/docs
      
      If you find Nessi useful, please consider starring the repository:
        https://github.com/nessi-dev/nessi
    EOS
  end
end
